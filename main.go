package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

var white_list_dirs []string

var white_list_dirs_nfs []string


var home_maker_temp_directory string
var uid_directory string


func get_user_from_id(id string) (user.User, error) {
	// get the filesystem type of the folder
	getent_out, err := exec.Command("getent", "passwd", id).Output()
	if err != nil {
		return user.User{}, err
	}

	passwd := strings.Split(string(getent_out), ":")
	gecos := strings.Split(passwd[4], ",")
	// log.Println(len(passwd))
	// log.Println(len(gecos))

	username := passwd[0]
	uid := passwd[2]
	gid := passwd[3]

	real_name := gecos[0]

	home_dir := passwd[5]
	// compare the filesystem type with "nfs"
	return user.User{Uid: uid, Gid: gid, Username: username, Name: real_name, HomeDir: home_dir}, nil
}



// isNFS returns true if the given folder is on an NFS mounted drive
func isNFS(folder string) (bool, error) {
	// get the filesystem type of the folder
	fsType, err := exec.Command("stat", "-f", "-c", "%T", folder).Output()
	if err != nil {
		return false, err
	}
	// compare the filesystem type with "nfs"
	return strings.TrimSpace(string(fsType)) == "nfs", nil
}


func mk_home_dir(dir string, uid int, gid int) (error) {
	will_make_dir := false
	base_dir := filepath.Join(dir, "..")

	if slices.Contains(white_list_dirs, base_dir) {
		will_make_dir = true
	}

	if slices.Contains(white_list_dirs_nfs, base_dir) {
		nfs_status, err := isNFS(base_dir)
		if err != nil {
			return err
		}

		if nfs_status {
			will_make_dir = true
		} else {
			will_make_dir = false
		}
	}



	if will_make_dir {
		os.Mkdir(dir, 0700)
		os.Chown(dir, uid, gid)
	}
	return nil
}



func make_dirs() {
	err := os.MkdirAll(home_maker_temp_directory, 0755)
	if err != nil {
		log.Fatal(err, " home_maker_temp_directory", " ", home_maker_temp_directory)
	}


	err = os.Chmod(home_maker_temp_directory, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(uid_directory, 0777)
	if err != nil {
		log.Fatal(err, " uid_directory")
	}


	err = os.Chmod(uid_directory, 0777)
	if err != nil {
		log.Fatal(err)
	}
}




func main() {
	home_maker_temp_directory = "/tmp/home-maker"
	uid_directory = home_maker_temp_directory+"/uid"

	make_dirs()

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// log.SetFlags(0) // removes all flags


	err := os.Chown(uid_directory, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	white_list_dirs = append(white_list_dirs, "/local-home")
	white_list_dirs_nfs = append(white_list_dirs_nfs, "/home")


	// defer dir.Close()

	for {
		// files, err := os.ReadDir(uid_directory)
		// if err != nil {
		//     log.Fatal(err)
		// }
		make_dirs()


		dir, err := os.Open(uid_directory)
		if err != nil {
			log.Fatal(err)
		}

	
		// Read the file names in the directory
		files, err := dir.Readdirnames(100) // Limit to the first 100 files
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
		}
	
		for _, file := range files {
			var user_cur_empty bool
			user_cur_empty = false

			user_cur, err := get_user_from_id(file)
			if err != nil {
				log.Println(err)
				user_cur_empty = true
			}

			var uid_int int
			var gid_int int

			if !user_cur_empty  {
				uid_int, err = strconv.Atoi(user_cur.Uid)
				if err != nil {
					log.Fatal(err)
				}
	
				gid_int, err = strconv.Atoi(user_cur.Gid)
				if err != nil {
					log.Fatal(err)
				}

				err = mk_home_dir(user_cur.HomeDir, uid_int, gid_int)
				if err != nil {
					log.Println(err)
				} else {
					// log.Println("made home folder for:", user_cur.Username, "/", user_cur.Name, "at:", user_cur.HomeDir)
					log.Printf("made home folder for: %s (%s) at: %s\n", user_cur.Name, user_cur.Username, user_cur.HomeDir)
				}
			}

			err = os.Remove(filepath.Join(uid_directory, file))

			if err != nil {
				log.Println(err)
			}
		}
		dir.Close()
		time.Sleep(3 * time.Second)

	}
}
