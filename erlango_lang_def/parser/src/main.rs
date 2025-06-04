use std::env;
use std::fs;
use std::process;

fn main() {
    // Collect command-line arguments
    let args: Vec<String> = env::args().collect();

    if args.len() < 3 {
        eprintln!("Usage: {} <bnf_file> <sample_input_file>", args[0]);
        process::exit(1);
    }

    let bnf_file_path = &args[1];
    let sample_file_path = &args[2];

    // Read BNF grammar file
    let bnf_content = match fs::read_to_string(bnf_file_path) {
        Ok(content) => content,
        Err(e) => {
            eprintln!("Error reading BNF file {}: {}", bnf_file_path, e);
            process::exit(1);
        }
    };

    println!("BNF grammar loaded from '{}'.", bnf_file_path);

    // Read sample input file
    let input_content = match fs::read_to_string(sample_file_path) {
        Ok(content) => content,
        Err(e) => {
            eprintln!("Error reading sample input file {}: {}", sample_file_path, e);
            process::exit(1);
        }
    };

    println!("\nSample input loaded from '{}':", sample_file_path);
    println!("----------------------------------------");
    println!("{}", input_content);
    println!("----------------------------------------");

    // Simple tokenizer: split by whitespace and common symbols
    println!("\nDetected tokens:");
    for line in input_content.lines() {
        let cleaned = line.trim();
        if cleaned.starts_with('#') || cleaned.is_empty() {
            continue; // Skip comment lines
        }

        let tokens = cleaned
            .split(|c: char| c.is_whitespace() || "(){}[],.".contains(c))
            .filter(|s| !s.is_empty());

        for token in tokens {
            println!("Token: {}", token);
        }
    }
}
