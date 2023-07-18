module github.com/pastor-robert/GoBox/mods/main

replace github.com/pastor-robert/GoBox/mods/m1 => ../m1

replace github.com/pastor-robert/GoBox/mods/m2 => ../m2

go 1.18

require (
	github.com/pastor-robert/GoBox/mods/m1 v0.0.0-00010101000000-000000000000 // indirect
	github.com/pastor-robert/GoBox/mods/m2 v0.0.0-00010101000000-000000000000 // indirect
)
