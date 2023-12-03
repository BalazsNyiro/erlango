# erlango
Erlang interpreter - written in Go

contact: Balazs Nyiro (balazs.nyiro.ca@gmail.com)

## Roadmap:

### in progress:
- [] Erlang source parse (raw token detection)                in progress
- [] Erlang code objects creation from tokens                 (2023 Jun)

### todo
- [] basic code execution                                     (2023 July)
- [] Erlang standard lib 1. implementation with version hooks (2023 Aug-Sep)
    - only the most important internal functions)

- [] Debugger tool building                                   (2023 Oct)
- [] speed optimization                                       (?)
- [] documentations, tutorials

- [] Erlang standard lib 2. implementation with version hooks (?) - this is a huge task, I hope

- [] signal sending/receiving with native Erlang instances (2025)
  (the first signal handling implementation will work with the Go version first)
  https://www.erlang.org/doc/reference_manual/processes.html#signals

## Why I write this interpreter?
I hope if:
- Go functions can be used natively from Erlang and vica-versa,
- the interpreter is quicker than the native Erlang version
- the debugger is far more useful

then this language can have a feature.

These two comments were the last drops in my glass.

### Is Erlang dead?
https://intellipaat.com/community/73770/is-erlang-dead
ErLang is one of the dying Programming Language
and the reason why its dying is because of difficulty in setting up
and very few developers who could support the Programming Language.

### Top 10 Programming Languages that will be Extinct in the year 2021
https://www.edureka.co/blog/top-10-dying-programming-languages
What is noteworthy is the fact that Erlang being a purely functional language
is not the sole factor behind Erlang’s decline. While there are still more jobs
for Erlang developers than there are developers available,
when it is compared to other languages, the demand for Erlang is a lot less.

### Go vs Rust
This is a difficult question. I know that Rust performance is better, and this was a serious hesitation for me - but:
 - Golang is concurrency focused (it was a big plus)
 - The human resource is extremely important. In Golang, the development speed is higher
 - I think the memory management and the compiler will be improved, so I hope Golang will be better with this.

## HALL OF FAME

In memory of Csanád Imreh (University of Szeged)

###  For the Erlang education:
- Melinda Tóth (Eötvös Loránd University)
- Zsolt Laky
- Robert Virding (special education in OTP Bank, Budapest)

###  For my friends in ERFI - we worked and enjoyed Erlang together:
- Ferenc Böröczki   (who was able to arrive with his bike in the coldest rain - and completely updated the team's daily codes during the nights :-)
- György Báló       (whose knowledge about the bank-system was incredible)
- Márkó Kitjonics   (who asked the best programming questions)
- Valentin Bujdosó  (who is a brilliant and kind man)
- Attila Faragó     (who was a good friend behind the hard surface)
- Zoltán Uramecz    (who was a Commodore 64 fan, like me)

- Péter Krekács     (who had the biggest wisdom and peace)
- Bence Szabó       (who always helped everybody - and took my bulb lamp at the end :-)
- Alex Boros        (who silently solved a tons of tasks)
- László Tóth       (we had a lot of discussions about functional programming)
- Balázs Boldóczki  (who supported our code in production)
- András Boroska    (who was my interviewer with Boldó, and gave me a chance to learn Erlang)
- Konyári Sándor    (who helped a lot to understand financial processes)
- László Popovics   (who himself believed in our toolset, and tested Erlang in special hardwares)
- Zoltán Dankó      (who had a vision, a dream about our team)


Thank you for our common time.

  


