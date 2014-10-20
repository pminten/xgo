package bob

var testCases = []struct {
	desc string
	in   string
	rep  int
	want string
}{
	{
		"stating something",
		"Tom-ay-to, tom-aaaah-to.",
		0,
		"Whatever.",
	},
	{
		"shouting",
		"WATCH OUT!",
		0,
		"Whoa, chill out!",
	},
	{
		"asking a question",
		"Does this cryogenic chamber make me look fat?",
		0,
		"Sure.",
	},
	{
		"asking a numeric question",
		"You are, what, like 15?",
		0,
		"Sure.",
	},
	{
		"talking forcefully",
		"Let's go make out behind the gym!",
		0,
		"Whatever.",
	},
	{
		"using acronyms in regular speech",
		"It's OK if you don't want to go to the DMV.",
		0,
		"Whatever.",
	},
	{
		"forceful questions",
		"WHAT THE HELL WERE YOU THINKING?",
		0,
		"Whoa, chill out!",
	},
	{
		"shouting numbers",
		"1, 2, 3 GO!",
		0,
		"Whoa, chill out!",
	},
	{
		"only numbers",
		"1, 2, 3",
		0,
		"Whatever.",
	},
	{
		"question with only numbers",
		"4?",
		0,
		"Sure.",
	},
	{
		"shouting with special characters",
		"ZOMG THE %^*@#$(*^ ZOMBIES ARE COMING!!11!!1!",
		0,
		"Whoa, chill out!",
	},
	{
		"shouting with no exclamation mark",
		"I HATE YOU",
		0,
		"Whoa, chill out!",
	},
	{
		"statement containing question mark",
		"Ending with ? means a question.",
		0,
		"Whatever.",
	},
	{
		"prattling on",
		"Wait! Hang on. Are you going to be OK?",
		0,
		"Sure.",
	},
	{
		"silence",
		"",
		0,
		"Fine. Be that way!",
	},
	{
		"prolonged silence",
		" ",
		10,
		"Fine. Be that way!",
	},
	{
		"alternate silences",
		"\t",
		10,
		"Fine. Be that way!",
	},
	{
		"multiple line questions",
		"Does this cryogenic chamber make me look fat?\nno",
		0,
		"Whatever.",
	},
}
