use std::{
    path::PathBuf,
    sync::mpsc,
    thread,
    time::Duration,
};

use clap::Parser;
use ratatui::{
    backend::CrosstermBackend,
    layout::Rect,
    Terminal,
};
use ratatui_image::{
    picker::Picker,
    protocol::StatefulProtocol,
    Image,
};

#[derive(Parser, Debug)]
struct Args {
    /// Previewer executable.
    #[arg(long, default_value = "")]
    previewer: String,

    /// File to preview.
    file: PathBuf,

    /// Preview width in terminal cells.
    width: u32,

    /// Preview height in terminal cells.
    height: u32,

    /// Horizontal position.
    x: u32,

    /// Vertical position.
    y: u32,

    /// "preview" or "preload".
    mode: String,
}

#[derive(Debug, Clone)]
struct PreviewArgs {
    filename: PathBuf,
    width: u32,
    height: u32,
    x: u32,
    y: u32,
    mode: String,
}

fn main() -> anyhow::Result<()> {
    let args = Args::parse();

    let preview_args = PreviewArgs {
        filename: args.file,
        width: args.width,
        height: args.height,
        x: args.x,
        y: args.y,
        mode: args.mode,
    };

    let (tx, rx) = mpsc::channel();

    thread::scope(|scope| {
        // Thread 1: load the image.
        let preview_thread = scope.spawn(|| {
            let picker = Picker::from_query_stdio().unwrap();

            let dyn_image = image::open(&preview_args.filename).unwrap();

            let protocol = picker
                .new_resize_protocol(dyn_image);

            tx.send(protocol).unwrap();
        });

        // Thread 2: do some other work.
        let test_thread = scope.spawn(|| {
            println!("Doing other work...");

            thread::sleep(Duration::from_secs(2));

            println!("Other work finished.");
        });

        // Receive the image from the preview thread.
        let mut protocol: StatefulProtocol = rx.recv().unwrap();

        // The terminal should be owned by this thread.
        let backend = CrosstermBackend::new(std::io::stdout());
        let mut terminal = Terminal::new(backend).unwrap();

        terminal.draw(|frame| {
            let area = Rect {
                x: preview_args.x,
                y: preview_args.y,
                width: preview_args.width,
                height: preview_args.height,
            };

            let image = Image::new(&mut protocol);
            frame.render_stateful_widget(image, area, &mut protocol);
        }).unwrap();

        // Wait for both worker threads.
        preview_thread.join().unwrap();
        test_thread.join().unwrap();
    });

    Ok(())
}
