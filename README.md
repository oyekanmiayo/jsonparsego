### Introduction
jsonparsergo A simple json parser I built live at the [SysDsgn Pop-Up](https://x.com/sysdsgn/status/1731710772714033305) on Dec 7th, 2023 in Lagos. I've made many modifications since and the final version can validate strings, numbers, literals (false, true and null), arrays and objects.

Some notes:
- I read the latest RFC for JSON ([RFC 8259](https://datatracker.ietf.org/doc/html/rfc8259)) which standardizes some inconsistencies from previous versions.
- JSON can represent four primitive types (strings, numbers, booleans, and null) and two structured types (objects and arrays)
- I first saw this challenge on [codingchallenges.fyi](https://codingchallenges.fyi/challenges/challenge-json-parser/)
  - I have tested against the recommended [suite](http://www.json.org/JSON_checker/); it passes all except fail25-fail28 because I think tabs and new lines in a string is fine.
  - I had to remove the backticks (`) from the json object in pass1 because it didn't play well with some Go syntax

### JSON Grammar
Here's the grammar that guided the validation checks in this project.

Structural Characters
```
begin-array     = ws %x5B ws  ; [ left square bracket

begin-object    = ws %x7B ws  ; { left curly bracket

end-array       = ws %x5D ws  ; ] right square bracket

end-object      = ws %x7D ws  ; } right curly bracket

name-separator  = ws %x3A ws  ; : colon

value-separator = ws %x2C ws  ; , comma

ws = *(
              %x20 /              ; Space
              %x09 /              ; Horizontal tab
              %x0A /              ; Line feed or New line
              %x0D )              ; Carriage return
```

Values 
```
value = false / null / true / object / array / number / string

false = %x66.61.6c.73.65   ; false

null  = %x6e.75.6c.6c      ; null

true  = %x74.72.75.65      ; true
```

Objects 

```
object = begin-object [ member *( value-separator member ) ]
               end-object

member = string name-separator value
```

Arrays 

```
array = begin-array [ value *( value-separator value ) ] end-array
```

Numbers

```
number = [ minus ] int [ frac ] [ exp ]

decimal-point = %x2E       ; .

digit1-9 = %x31-39         ; 1-9

e = %x65 / %x45            ; e E

exp = e [ minus / plus ] 1*DIGIT

frac = decimal-point 1*DIGIT

int = zero / ( digit1-9 *DIGIT )

minus = %x2D               ; -

plus = %x2B                ; +

zero = %x30                ; 0
```

Strings
```
string = quotation-mark *char quotation-mark

char = unescaped /
  escape (
      %x22 /          ; "    quotation mark  U+0022
      %x5C /          ; \    reverse solidus U+005C
      %x2F /          ; /    solidus         U+002F
      %x62 /          ; b    backspace       U+0008
      %x66 /          ; f    form feed       U+000C
      %x6E /          ; n    line feed       U+000A
      %x72 /          ; r    carriage return U+000D
      %x74 /          ; t    tab             U+0009
      %x75 4HEXDIG )  ; uXXXX                U+XXXX

escape = %x5C              ; \

quotation-mark = %x22      ; "

unescaped = %x20-21 / %x23-5B / %x5D-10FFFF
```