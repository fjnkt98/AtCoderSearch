package list

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fmt"
	"regexp"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
)

type Language struct {
	Language string `db:"language"`
	Group    string `db:"group"`
}

var MajorLanguageMapping = map[string]string{
	"Ada2012 (GNAT 9.2.1)":                         "Other",
	"Awk (GNU Awk 4.1.4)":                          "Awk",
	"AWK (GNU Awk 5.0.1)":                          "Awk",
	"Awk (mawk 1.3.3)":                             "Awk",
	"Bash (4.2.25)":                                "Bash",
	"Bash (5.0.11)":                                "Bash",
	"Bash (bash 5.2.2)":                            "Bash",
	"Bash (GNU bash v4.3.11)":                      "Bash",
	"bc (1.07.1)":                                  "Other",
	"Brainfuck (bf 20041219)":                      "Other",
	"C# 11.0 AOT (.NET 7.0.7)":                     "C#",
	"C# 11.0 (.NET 7.0.7)":                         "C#",
	"C++11 (Clang++ 3.4)":                          "C++",
	"C++11 (GCC 4.8.1)":                            "C++",
	"C++11 (GCC 4.9.2)":                            "C++",
	"C++14 (Clang++ 3.4)":                          "C++",
	"C++14 (Clang 3.8.0)":                          "C++",
	"C++14 (GCC 5.4.1)":                            "C++",
	"C++ 17 (Clang 16.0.5)":                        "C++",
	"C++ 17 (gcc 12.2)":                            "C++",
	"C++ 20 (Clang 16.0.5)":                        "C++",
	"C++ 20 (gcc 12.2)":                            "C++",
	"C++ 23 (Clang 16.0.5)":                        "C++",
	"C++ 23 (gcc 12.2)":                            "C++",
	"Carp (Carp 0.5.5)":                            "Other",
	"C (Clang 10.0.0)":                             "C",
	"C++ (Clang 10.0.0)":                           "C++",
	"C++ (Clang 10.0.0 with AC Library)":           "C++",
	"C++ (Clang 10.0.0 with AC Library v1.0)":      "C++",
	"C++ (Clang 10.0.0 with AC Library v1.1)":      "C++",
	"C (Clang 3.4)":                                "C",
	"C++ (Clang++ 3.4)":                            "C++",
	"C (Clang 3.8.0)":                              "C",
	"C++ (Clang 3.8.0)":                            "C++",
	"Ceylon (1.2.1)":                               "Other",
	"C++ (G++ 4.6.4)":                              "C++",
	"C (gcc 12.2.0)":                               "C",
	"C (GCC 4.4.7)":                                "C",
	"C++ (GCC 4.4.7)":                              "C++",
	"C (GCC 4.6.4)":                                "C",
	"C (GCC 4.9.2)":                                "C",
	"C++ (GCC 4.9.2)":                              "C++",
	"C (GCC 5.4.1)":                                "C",
	"C++ (GCC 5.4.1)":                              "C++",
	"C (GCC 9.2.1)":                                "C",
	"C++ (GCC 9.2.1)":                              "C++",
	"C++ (GCC 9.2.1 with AC Library)":              "C++",
	"C++ (GCC 9.2.1 with AC Library v1.0)":         "C++",
	"C++ (GCC 9.2.1 with AC Library v1.1)":         "C++",
	"Clojure (1.10.1.536)":                         "Clojure",
	"Clojure (1.1.0 + OpenJDK 1.7)":                "Clojure",
	"Clojure (1.8.0)":                              "Clojure",
	"Clojure (babashka 1.3.181)":                   "Clojure",
	"Clojure (clojure 1.11.1)":                     "Clojure",
	"C# (Mono 2.10.8.1)":                           "C#",
	"C# (Mono 3.2.1.0)":                            "C#",
	"C# (Mono 4.6.2.0)":                            "C#",
	"C# (Mono-csc 3.5.0)":                          "C#",
	"C# (Mono-mcs 6.8.0.105)":                      "C#",
	"C# (.NET Core 3.1.201)":                       "C#",
	"COBOL - Fixed (OpenCOBOL 1.1.0)":              "COBOL",
	"COBOL (Free) (GnuCOBOL 3.1.2)":                "COBOL",
	"COBOL - Free (OpenCOBOL 1.1.0)":               "COBOL",
	"COBOL (GnuCOBOL(Fixed) 3.1.2)":                "COBOL",
	"Common Lisp (SBCL 1.0.55.0)":                  "Common Lisp",
	"Common Lisp (SBCL 1.1.14)":                    "Common Lisp",
	"Common Lisp (SBCL 2.0.3)":                     "Common Lisp",
	"Common Lisp (SBCL 2.3.6)":                     "Common Lisp",
	"Crystal (0.20.5)":                             "Ruby",
	"Crystal (0.33.0)":                             "Ruby",
	"Crystal (Crystal 1.9.1)":                      "Ruby",
	"Cyber (Cyber v0.2-Latest)":                    "Other",
	"Cython (0.29.16)":                             "Python",
	"Dart (2.7.2)":                                 "Dart",
	"Dart (Dart 3.0.5)":                            "Dart",
	"Dash (0.5.8)":                                 "Other",
	"dc (1.4.1)":                                   "Other",
	"D (DMD 2.060)":                                "D",
	"D (DMD 2.066.1)":                              "D",
	"D (DMD 2.091.0)":                              "D",
	"D (DMD 2.104.0)":                              "D",
	"D (DMD64 v2.070.1)":                           "D",
	"D (GDC 4.9.4)":                                "D",
	"D (GDC 9.2.1)":                                "D",
	"D (LDC 0.17.0)":                               "D",
	"D (LDC 1.20.1)":                               "D",
	"D (LDC 1.32.2)":                               "D",
	"ECLiPSe (ECLiPSe 7.1_13)":                     "Other",
	"Elixir (1.10.2)":                              "Elixir",
	"Elixir (Elixir 1.15.2)":                       "Elixir",
	"Emacs Lisp (Native Compile) (GNU Emacs 28.2)": "Other",
	"Erlang (22.3)":                                "Erlang",
	"F# 7.0 (.NET 7.0.7)":                          "F#",
	"><> (fishr 0.1.0)":                            "Other",
	"F# (Mono 10.2.3)":                             "F#",
	"F# (Mono 4.0)":                                "F#",
	"F# (.NET Core 3.1.201)":                       "F#",
	"Forth (gforth 0.7.3)":                         "Other",
	"Fortran (gfortran 12.2)":                      "Fortran",
	"Fortran (gfortran v4.8.4)":                    "Fortran",
	"Fortran (GNU Fortran 9.2.1)":                  "Fortran",
	"Fortran(GNU Fortran 9.2.1)":                   "Fortran",
	"Go (1.14.1)":                                  "Go",
	"Go (1.4.1)":                                   "Go",
	"Go (1.6)":                                     "Go",
	"Go (go 1.20.6)":                               "Go",
	"Haskell (GHC 7.10.3)":                         "Haskell",
	"Haskell (GHC 7.4.1)":                          "Haskell",
	"Haskell (GHC 8.8.3)":                          "Haskell",
	"Haskell (GHC 9.4.5)":                          "Haskell",
	"Haskell (Haskell Platform 2014.2.0.0)":        "Haskell",
	"Haxe (4.0.3); Java":                           "Other",
	"Haxe (4.0.3); js":                             "Other",
	"IOI-Style C++ (GCC 5.4.1)":                    "Other",
	"Java7 (OpenJDK 1.7.0)":                        "Java",
	"Java8 (OpenJDK 1.8.0)":                        "Java",
	"Java (OpenJDK 11.0.6)":                        "Java",
	"Java (OpenJDK 17)":                            "Java",
	"Java (OpenJDK 1.7.0)":                         "Java",
	"Java (OpenJDK 1.8.0)":                         "Java",
	"JavaScript (Deno 1.35.1)":                     "JavaScript",
	"JavaScript (Node.js 0.6.12)":                  "JavaScript",
	"JavaScript (Node.js 12.16.1)":                 "JavaScript",
	"JavaScript (Node.js 18.16.1)":                 "JavaScript",
	"JavaScript (Node.js v0.10.36)":                "JavaScript",
	"JavaScript (node.js v5.12)":                   "JavaScript",
	"jq (jq 1.6)":                                  "Other",
	"Julia (0.5.0)":                                "Julia",
	"Julia (1.4.0)":                                "Julia",
	"Julia (Julia 1.9.2)":                          "Julia",
	"Koka (koka 2.4.0)":                            "Other",
	"Kotlin (1.0.0)":                               "Kotlin",
	"Kotlin (1.3.71)":                              "Kotlin",
	"Kotlin (Kotlin/JVM 1.8.20)":                   "Kotlin",
	"LLVM IR (Clang 16.0.5)":                       "Other",
	"Lua (5.3.2)":                                  "Lua",
	"LuaJIT (2.0.4)":                               "Lua",
	"Lua (Lua 5.3.5)":                              "Lua",
	"Lua (Lua 5.4.6)":                              "Lua",
	"Lua (LuaJIT 2.1.0)":                           "Lua",
	"Lua (LuaJIT 2.1.0-beta3)":                     "Lua",
	"MoonScript (0.5.0)":                           "Other",
	"Nibbles (literate form) (nibbles 1.01)":       "Other",
	"Nim (0.13.0)":                                 "Nim",
	"Nim (1.0.6)":                                  "Nim",
	"Nim (Nim 1.6.14)":                             "Nim",
	"Objective-C (Clang 10.0.0)":                   "Objective-C",
	"Objective-C (Clang3.8.0)":                     "Objective-C",
	"Objective-C (GCC 5.3.0)":                      "Objective-C",
	"OCaml (3.12.1)":                               "Ocaml",
	"OCaml (4.02.1)":                               "Ocaml",
	"OCaml (4.02.3)":                               "Ocaml",
	"OCaml (4.10.0)":                               "Ocaml",
	"OCaml (ocamlopt 5.0.0)":                       "Ocaml",
	"Octave (4.0.2)":                               "Octave",
	"Octave (5.2.0)":                               "Octave",
	"Octave (GNU Octave 8.2.0)":                    "Octave",
	"Pascal (fpc 2.4.4)":                           "Pascal",
	"Pascal (FPC 2.6.2)":                           "Pascal",
	"Pascal (FPC 3.0.4)":                           "Pascal",
	"Pascal (fpc 3.2.2)":                           "Pascal",
	"Perl (5.14.2)":                                "Perl",
	"Perl (5.26.1)":                                "Perl",
	"Perl6 (rakudo-star 2016.01)":                  "Perl",
	"Perl (perl  5.34)":                            "Perl",
	"Perl (v5.18.2)":                               "Perl",
	"PHP (5.6.30)":                                 "PHP",
	"PHP (7.4.4)":                                  "PHP",
	"PHP7 (7.0.15)":                                "PHP",
	"PHP (PHP 5.3.10)":                             "PHP",
	"PHP (PHP 5.5.21)":                             "PHP",
	"PHP (php 8.2.8)":                              "PHP",
	"PowerShell (PowerShell 7.3.1)":                "PowerShell",
	"Prolog (SWI-Prolog 8.0.3)":                    "Prolog",
	"Prolog (SWI-Prolog 9.0.4)":                    "Prolog",
	"PyPy2 (5.6.0)":                                "Python",
	"PyPy2 (7.3.0)":                                "Python",
	"PyPy3 (2.4.0)":                                "Python",
	"PyPy3 (7.3.0)":                                "Python",
	"Python2 (2.7.6)":                              "Python",
	"Python (2.7.3)":                               "Python",
	"Python (2.7.6)":                               "Python",
	"Python3 (3.2.3)":                              "Python",
	"Python3 (3.4.2)":                              "Python",
	"Python3 (3.4.3)":                              "Python",
	"Python (3.4.3)":                               "Python",
	"Python (3.8.2)":                               "Python",
	"Python (CPython 3.11.4)":                      "Python",
	"Python (Cython 0.29.34)":                      "Python",
	"Python (Mambaforge / CPython 3.10.10)":        "Python",
	"Python (PyPy 3.10-v7.3.12)":                   "Python",
	"Racket (7.6)":                                 "Other",
	"Raku (Rakudo 2020.02.1)":                      "Other",
	"Raku (Rakudo 2023.06)":                        "Other",
	"ReasonML (reason 3.9.0)":                      "Other",
	"R (GNU R 4.2.1)":                              "R",
	"Ruby (1.9.3)":                                 "Ruby",
	"Ruby (1.9.3p550)":                             "Ruby",
	"Ruby (2.1.5p273)":                             "Ruby",
	"Ruby (2.3.3)":                                 "Ruby",
	"Ruby (2.7.1)":                                 "Ruby",
	"Ruby (ruby 3.2.2)":                            "Ruby",
	"Rust (1.15.1)":                                "Rust",
	"Rust (1.42.0)":                                "Rust",
	"Rust (rustc 1.70.0)":                          "Rust",
	"SageMath (SageMath 9.5)":                      "Other",
	"Scala (2.11.5)":                               "Scala",
	"Scala (2.11.7)":                               "Scala",
	"Scala (2.13.1)":                               "Scala",
	"Scala (2.9.1)":                                "Scala",
	"Scala 3.3.0 (Scala Native 0.4.14)":            "Scala",
	"Scala (Dotty 3.3.0)":                          "Scala",
	"Scheme (Gauche 0.9.1)":                        "Scheme",
	"Scheme (Gauche 0.9.12)":                       "Scheme",
	"Scheme (Gauche 0.9.3.3)":                      "Scheme",
	"Scheme (Gauche 0.9.9)":                        "Scheme",
	"Scheme (Scheme 9.1)":                          "Scheme",
	"Sed (4.4)":                                    "Sed",
	"Sed (GNU sed 4.2.2)":                          "Sed",
	"Sed (GNU sed 4.8)":                            "Sed",
	"Seed7 (Seed7 3.2.1)":                          "Other",
	"Standard ML (MLton 20100608)":                 "Other",
	"Standard ML (MLton 20130715)":                 "Other",
	"Swift (5.2.1)":                                "Swift",
	"Swift (swift-2.2-RELEASE)":                    "Swift",
	"Swift (swift 5.8.1)":                          "Swift",
	"Text (cat)":                                   "Text",
	"Text (cat 8.28)":                              "Text",
	"Text (cat 8.32)":                              "Text",
	"TypeScript (2.1.6)":                           "TypeScript",
	"TypeScript (3.8)":                             "TypeScript",
	"TypeScript 5.1 (Deno 1.35.1)":                 "TypeScript",
	"TypeScript 5.1 (Node.js 18.16.1)":             "TypeScript",
	"Unlambda (0.1.3)":                             "Other",
	"Unlambda (2.0.0)":                             "Other",
	"Vim (8.2.0460)":                               "Other",
	"Visual Basic (Mono 2.10.8)":                   "Visual Basic",
	"Visual Basic (Mono 4.0.1)":                    "Visual Basic",
	"Visual Basic (.NET Core 3.1.101)":             "Visual Basic",
	"V (V 0.4)":                                    "Other",
	"Whitespace (whitespacers 1.0.0)":              "Other",
	"Zig (Zig 0.10.1)":                             "Zig",
	"Zsh (5.4.2)":                                  "Other",
	"なでしこ (cnako3 3.4.20)":                         "Other",
	"プロデル (mono版プロデル 1.9.1182)":                    "Other",
}

var LanguagePatterns = map[string]*regexp.Regexp{
	"Awk":          regexp.MustCompile(`(?i)^awk`),
	"Bash":         regexp.MustCompile(`(?i)^bash`),
	"C":            regexp.MustCompile(`(?i)^C `),
	"C#":           regexp.MustCompile(`(?i)^C#`),
	"C++":          regexp.MustCompile(`(?i)^C\+\+`),
	"COBOL":        regexp.MustCompile(`(?i)^COBOL`),
	"Clojure":      regexp.MustCompile(`(?i)^Clojure`),
	"Common Lisp":  regexp.MustCompile(`(?i)^Common Lisp`),
	"D":            regexp.MustCompile(`(?i)^D `),
	"Dart":         regexp.MustCompile(`(?i)^Dart`),
	"Elixir":       regexp.MustCompile(`(?i)^Elixir`),
	"Erlang":       regexp.MustCompile(`(?i)^Erlang`),
	"F#":           regexp.MustCompile(`(?i)^F#`),
	"Fortran":      regexp.MustCompile(`(?i)^Fortran`),
	"Go":           regexp.MustCompile(`(?i)^Go`),
	"Haskell":      regexp.MustCompile(`(?i)^Haskell`),
	"Java":         regexp.MustCompile(`(?i)^Java\d{0,2} `),
	"JavaScript":   regexp.MustCompile(`(?i)^JavaScript`),
	"Julia":        regexp.MustCompile(`(?i)^Julia`),
	"Kotlin":       regexp.MustCompile(`(?i)^Kotlin`),
	"Lua":          regexp.MustCompile(`(?i)^Lua`),
	"Nim":          regexp.MustCompile(`(?i)^Nim`),
	"Objective-C":  regexp.MustCompile(`(?i)^Objective-C`),
	"Ocaml":        regexp.MustCompile(`(?i)^Ocaml`),
	"Octave":       regexp.MustCompile(`(?i)^Octave`),
	"PHP":          regexp.MustCompile(`(?i)^PHP`),
	"Pascal":       regexp.MustCompile(`(?i)^Pascal`),
	"Perl":         regexp.MustCompile(`(?i)^Perl`),
	"PowerShell":   regexp.MustCompile(`(?i)^PowerShell`),
	"Prolog":       regexp.MustCompile(`(?i)^Prolog`),
	"Python":       regexp.MustCompile(`(?i)^([PC]ython|PyPy)`),
	"R":            regexp.MustCompile(`(?i)^R `),
	"Ruby":         regexp.MustCompile(`(?i)^(Ruby|Crystal)`),
	"Rust":         regexp.MustCompile(`(?i)^Rust`),
	"Scala":        regexp.MustCompile(`(?i)^Scala`),
	"Scheme":       regexp.MustCompile(`(?i)^Scheme`),
	"Sed":          regexp.MustCompile(`(?i)^Sed`),
	"Swift":        regexp.MustCompile(`(?i)^Swift`),
	"Text":         regexp.MustCompile(`(?i)^Text`),
	"TypeScript":   regexp.MustCompile(`(?i)^TypeScript`),
	"Visual Basic": regexp.MustCompile(`(?i)^Visual Basic`),
	"Zig":          regexp.MustCompile(`(?i)^Zig`),
}

func GetLanguageGroup(language string) string {
	if group, ok := MajorLanguageMapping[language]; ok {
		return group
	}

	for group, pattern := range LanguagePatterns {
		if pattern.MatchString(language) {
			return group
		}
	}

	return "Other"
}

func Update(ctx context.Context, db *sqlx.DB) error {
	rows, err := db.QueryContext(
		ctx,
		`
		SELECT
			DISTINCT "language"
		FROM
			"submissions"	
		`,
	)
	if err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to fetch languages from submissions table"))
	}
	defer rows.Close()

	languages := make([]Language, 0, 32)
	for rows.Next() {
		var lang string
		if err := rows.Scan(&lang); err != nil {
			return failure.Translate(err, acs.DBError, failure.Message("failed to scan row"))
		}
		group := GetLanguageGroup(lang)
		languages = append(languages, Language{Language: lang, Group: group})
	}

	tx, err := db.Beginx()
	if err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to start transaction"))
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx,
		`
		DELETE FROM "languages"
		`,
	); err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to delete languages"))
	}

	affected := 0
	if result, err := tx.NamedExecContext(
		ctx,
		`
		INSERT INTO "languages" ("language", "group") VALUES (:language, :group)
		`,
		languages,
	); err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to insert languages"))
	} else {
		a, _ := result.RowsAffected()
		affected = int(a)
	}

	if err := tx.Commit(); err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to commit transaction to save languages"))
	} else {
		slog.Info(fmt.Sprintf("commit transaction. save %d rows.", affected))
	}

	return nil
}
