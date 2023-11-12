# The Linux Programming Interface

This repository contains the source code for my solutions to the programming exercises in [_The Linux Programming Interface_](https://man7.org/tlpi/) book, by Michael Kerrisk.

> [!NOTE]
>
> The solutions are written in `GO` on `Arch Linux` with kernel version `6.5.4-arch2-1`.

## Repository structure

The book is divided into 64 chapters, some of which contains a set of exercises at the end. Each chapter has its own directory, named `chXX`, where `XX` is the chapter number. Each of these directories contains a `Makefile` that can be used to build the solutions for that chapter.

The repository is structured as follows:

```bash
chXX
├── bin
├── makefile
├── README.md
└── src
    ├── <exercise_name>.go
    └── ...
```

## Requirements

- `go` compiler (version `1.21.3` or higher)
- `make` (version `4.4.1` or higher)

## Local setup

To setup the repository locally, clone it and run `make` in the root directory:

```bash
git clone git@github.com:ahmeducf/tlpi.git
cd tlpi
make
```

## Building the solutions

To build the solutions for a chapter, run `make` in the corresponding directory. For example, to build the solutions for chapter 4, run:

```bash
cd ch04
make
```

## Running the solutions

The solutions are built as executables in `bin` directory in the same chapter directory with the same name as the source file. For example, the solution for `tee` program in ch04 is built as `tee`. To run it, simply execute it:

```bash
ch04/bin/tee
```

## License

The source code is licensed under the [MIT License](https://opensource.org/licenses/MIT).
