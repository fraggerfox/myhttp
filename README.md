# URL content MD5 hasher

This repository contains the URL content MD5 hasher exercise given as a part of
the [adjust's][1] recruitment process.

The implementation is done using [Golang][2].

## Implementation details

 - The code is written in a procedural fashion.
 - Unit tests are written for the appropriate sections of the code.
 - The only external library used is the extension for Golang's testing library
   called [check.v1][3]

 ## How to build the program

 - The repository comes with `Makefile` which has appropriate targets to build /
   test the program.
 - The `Makefile` has a `help` target that briefly explains what each of the
   target does.

## How to run the program

- The program is a command line program and the user is expected to run it in
  a command line environment by invoking the executable name.
- The `-parallel` flags allows to set the maximum number of parallel instances
  to run when fetching URL.
- Example outputs
  - Invoking without any parameters
    ```
	$ ./myhttp
	Usage: myhttp [-parallel JOBS] URL1 URL2...

	Checks if a give URL is valid and generates the MD5 of the contents,
	Example: myhttp -parallel 1 http://adjust.com

	Interpretation of parameters:
        JOBS            Number simultaneous jobs to process the URLs. Defaults to 10.
        URL             A list of space separated URLs to fetch content and generate MD5.
                        For example adjust.com http://google.com
    ```
  - Invoking with parameters
    ```
    $ ./myhttp www.google.com http://adjust.com
	http://www.google.com: d6207fd31973352782940d5aff6c6b34
	http://adjust.com: 7ccbb8a0f7bb1f6f009c4ccd2184c1cb
    $
    ```
    ```
    $ ./myhttp -parallel 5 adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com
	http://adjust.com: 7ccbb8a0f7bb1f6f009c4ccd2184c1cb
	http://yandex.com: 88bde7b9bded3be39a416575f176afb7
	http://google.com: 0ff200fe2ec0b2e00bff6d7c3dee45fb
	http://twitter.com: 19ff5073884551896b867e677373cb5e
	http://yahoo.com: b770521149e477855da3777a9d23c2bc
	http://baroquemusiclibrary.com: 24d95e220af341299f1a1f1553b7ea60
	http://reddit.com/r/notfunny: 88a4ef066fb25af688ba77fc06590d01
	http://reddit.com/r/funny: 01ad2392ac99af48dfa7db1976f85226
	http://facebook.com: 44b9acf74e0eaee607a0c9ce8dc08e7d
    $
    ```

## License

BSD 2-clause. See [LICENSE][4].

[1]: https://www.stackbuilders.com/
[2]: https://golang.org
[3]: https://gopkg.in/check.v1
[4]: https://github.com/fraggerfox/myhttp/blob/master/LICENSE
