# Shamir: A Command-Line Extraction of Vault's Shamir's Secret Sharing Implementation

## WARNING: DO NOT USE THIS

This is probably a bad idea and poorly conceived. There's a strong likelihood
that I will shoot myself in the foot with this (as will you if you use it).
This is me rolling my own crypto.

## Background

I want my 1Password password/secret key to be accessible by my family in case
something happens to me (eaten by a bear, bonk my head skiing, forswear
technology and join a monastic order). I don't, however, want it to be
accessible unilaterally or accessible to a burglar who just happens to find the
Post-It note.

I want the secrets to be split up, such that the secrets can only be recovered
if some number of splits are put together. Luckily, [Shamir's Secret
Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) does exactly
this. Even more fortuitously, there is a golang [implementation in Hashicorp
Vault](https://github.com/hashicorp/vault/tree/master/shamir). Serendipitously,
I also found myself in between jobs and wanting to play with this concept.

With this, I will split my password and secret key and distribute them to
trusted parties (who are also not too, _too_ close to each other, to avoid
collusion).

## Set up

This requires setting up some developer tools to run the program, this might
get a little frustrating (get one of my nerd friends to do it).

_TODO: Add Docker instructions_

1. [Install Go](https://golang.org/doc/install). Follow the instructions there
   for your operating system. (Go is the programming language this tool is
   written in. Installing Go will allow you to build and run programs written
   in the Go language, like this one)
2. Download [this repo as a Zip
   file](https://github.com/schleyfox/shamir/archive/main.zip) and unzip it
   (probably by double-clicking on it). Make note of the file path (on Mac this
   should be `~/Downloads/shamir`)
2. Open a Terminal.app window (Mac) or Command Prompt (Windows) (I assume Linux
   users already have their preferred term open and are using it to try to fix
   their Wi-Fi or printer)
3. In the terminal, change the directory to the path you unzipped this code to.
   For example, `cd ~/Downloads/shamir`
4. Build the project, with `go build -o ./bin/shamir`. Keep the terminal open
   for the next section.

## Usage

### Combine (retrieving secrets from the splits)

1. In a terminal at the path where you built the project (e.g.
   `~/Downloads/shamir`), run `bin/shamir combine --parts=<PARTS>` where
   `<PARTS>` is the number of splits you have access to, like `bin/shamir
   combine --parts=2` (note that this number needs to be greater than the
   threshold required, otherwise you'll have insufficient information to
   recover the secrets)
2. Enter each split you have when prompted. There is some error detection to
   guard against typos. If errors are detected, re-enter the split correctly.
3. Once `<PARTS>` splits have been entered, the secret will be output

Example:

```
$ bin/shamir combine --parts=2
Share 0: /n1KUOPtx0DsMg==
Share 1: vIIqNd6I8VcMoA==
# SECRET
atest
```

Note that we do not turn off tty echo for share entries, as we are outputting
the re-formed secret anyway. Additionally, it would be frustrating and error
prone to not be able to see the random values you were typing in.

Also, there's no protection built in against people providing fake shares. The
checksum only validates the integrity of the share by itself, not that it is
the correct share for this secret. If it's wrong, the output secret will be
wrong (most likely will be quite obviously wrong)

### Split (split a secret into splits/shares/parts)

1. In a terminal at the path where you built the project (e.g. `~/Downloads/shamir`) run `bin/shamir split --parts=<PARTS> --threshold=<THRESHOLD>` where `<PARTS>` is the number of shares you wish to produce and `<THRESHOLD>` is the number of those shares that need to be brought together to reform the secret. For example if you want to give 7 people a share and require that 4 of them be present to unlock the secret, you could run `bin/shamir split --parts=7 --threshold=4`.
2. Enter the secret when prompted.
3. The splits will be output along with the command to recombine them.

Example:

```
$ bin/shamir split --parts=3 --threshold=2           20:40:57
Secret: atest
# parts = 3, threshold = 2
# bin/shamir combine --parts=2
/n1KUOPtx0DsMg==
G+MVC52aUFFt1Q==
vIIqNd6I8VcMoA==
```

Note that we do not turn off tty echo for the secret entry to allow errors in
entry to be checked. Also, I didn't feel like messing with tty stuff for this
project.

## Test Splits

```
# parts = 4, threshold = 3
# bin/shamir combine --parts=3
1K/DtEEgGaubj3UQwGA=
alqED3Z76yCKtYpmniw=
q9OVJY7oMCevvdu9zFA=
d+1p9BnbqxxveYQHgS0=
```

```
# parts = 3, threshold = 2
# bin/shamir combine --parts=2
LEWUi4Iu6p7aSeH9D/USIbtL6G0=
PDmxSJbrPkD8hUEXt1aVL19x9dQ=
XwcznUeIc+hsVpNU+9HkBjWzlKs=
```
