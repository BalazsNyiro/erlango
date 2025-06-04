use std::fs::File;
use std::io::{BufRead, BufReader};
use regex::Regex;

fn main() -> std::io::Result<()> {
    // Load BNF file
    let _grammar = std::fs::read_to_string("erlango_lang_balazs.bnf")?;
    println!("BNF grammar loaded.");

    // Define regex patterns for atom and number based on BNF
    let atom_re = Regex::new(r"^[a-z][a-zA-Z0-9_@]*$").unwrap();
    let number_re = Regex::new(r"^\d+(\.\d+)?$").unwrap();

    // Load and parse sample input
    let file = File::open("sample_input.txt")?;
    let reader = BufReader::new(file);

    for (i, line) in reader.lines().enumerate() {
        let line = line?.trim().to_string();
        if line.is_empty() {
            continue;
        }

        if atom_re.is_match(&line) {
            println!("Line {}: '{}' is a valid ATOM", i + 1, line);
        } else if number_re.is_match(&line) {
            println!("Line {}: '{}' is a valid NUMBER", i + 1, line);
        } else {
            println!("Line {}: '{}' is INVALID", i + 1, line);
        }
    }

    Ok(())
}