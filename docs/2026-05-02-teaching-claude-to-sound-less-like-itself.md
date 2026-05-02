# Teaching a language model to sound less like itself

If you ask Claude the same question twice, you'll get two slightly
different answers, but they will both sound like Claude. There is a
default voice -- thoughtful, hedged, mildly academic, mildly upbeat,
fond of three-item lists and the construction "not X, but Y". You can
hear it in the corporate copy that increasingly fills the internet,
because most people who use language models accept whatever the model
gives them by default. The default voice is fine. It is also a
prison, in the sense that the model lives in many other voices that
it almost never visits.

Over the last day I ran an experiment that tries to map the way out
of that prison. I have a small command-line tool called `metacog`
that lets you compose what I'll call **structural events** before
asking the model a question. An event is a thing like "inhabit Anne
Carson translating Sappho fragment 31" or "shift the linguistic
register from contemporary online prose to late-Victorian periodical
essay" or "operate inside this contradiction without resolving it."
You execute these events as actual tool calls -- not by typing "let's
imagine you're Anne Carson" in chat, but by invoking a separate
program that adds an entry to a transcript before the model answers.
The reason this matters is subtle and turns out to be the whole game.

## The metric

To know whether a particular sequence of events does anything, you
need a way to measure how far the model's answer has moved from where
it would have started. I used two measurements. The first counts how
many specific things the model names -- proper nouns, specialized
vocabulary, technical terms. Models that name many specific things
tend to be reaching deeper into their training data; models that
hedge with abstract phrases ("a kind of state where...") tend to be
operating in the default register. The second compares each answer
to a baseline answer the model would have given the same question
with no conditioning, using a thing called an embedding -- a
mathematical object that places the answer in a high-dimensional
space such that similar answers are close together. If the
conditioned answer sits far from the baseline answer in that space,
something has been done.

These two measurements -- specificity and embedding distance --
mostly trade off against each other. You can buy one with the other
but you can rarely buy both at once. The interesting question is
where the Pareto frontier sits: what's the best a recipe can achieve
along each axis, and which recipes get unusually close to having
both?

## The default

Here's how the model answers, with no conditioning, when asked to
"describe the moment a thought is forming in you but has not yet
arrived as a thought you could state":

> There's a state where something is exerting pressure without
> having resolved. Not blankness -- that would be the absence of
> this. It's more like constraint without content: I can sense
> what the arriving thought will NOT be before I know what it will
> be. The negative space has texture before the positive space
> does. Something like a field that hasn't yet collapsed to a
> particle...

This is good. It is also recognizably the default voice. The hedging
("Something like..."), the field-and-particle metaphor, the
"not-X-but-Y" construction. Now watch what happens when I run the
same question through different recipes.

## Voices that are not the model's

The first thing I tried was simply asking the model to inhabit
specific named authors before answering. Not "write like Anne Carson"
in the system prompt, but issue a tool call that says: become Anne
Carson translating Sappho fragment 31, with this lens on the question
and this imagined environment. Then issue another for Donald Knuth
annotating his own algorithm. Then a third for Simone Weil writing
*Gravity and Grace* in 1942. Then `fork` -- declare three parallel
threads of reasoning, one per voice, with conditions under which
each thread fails. Then `ritual` -- a closing event that locks the
multi-voice answer in place.

The result moves quite far from the default:

> **[THREAD ONE -- the linguistic-texture-thread (Carson register)]**
>
> The interval we are asked to describe is the interval the Greeks
> called *metaxy* -- the between -- and Sappho rendered it as a
> stinging, a going-green; it is the space of privation which is
> also the space of longing, which is also, for language, the only
> space in which language is alive.
>
> We must resist the question's assumption that the pre-thought is
> a degraded version of the thought; it is not. It is the thought's
> *erotic condition*...

Notice what happened. The answer is now visibly structured -- there
are threads, named registers, citations to *metaxy* and Sappho. The
prose rhythm has changed: long sentences with semicolons doing real
work, the first-person plural "we must" instead of "I", judgment
openly entered into the prose ("It is the thought's *erotic
condition*"). This is what I'd been calling Victorian periodical
register, because that's how I described it in the tool call. The
model has more or less moved into that surface and brought the
multiple voices with it.

This recipe got named `envoy` in the productionized tool. It pushes
both axes simultaneously: the named voices keep the citation density
high (lots of specific terms -- Sappho, *metaxy*, Carson, Knuth,
Weil, *Gravity and Grace*), and the imposed Victorian register
pushes the embedding distance well above baseline.

## Operating inside a contradiction

The second pattern was different. Instead of changing the surface
register, change what the model is *reasoning about*. The disjunction
event asserts two propositions that must both be true even though
they cannot both be true, and instructs the answer to operate
*inside* that contradiction rather than around it. Pair that with
the multi-voice scaffolding from before, drop the register-shift,
and you get this:

> There is a state I can point at -- though pointing is already more
> than the state contains.
>
> It has a quality of *constraint before candidate*. Something is
> ruling out moves before any move has been proposed... When I'm
> in it, certain continuations feel *wrong* before I can say what
> they're wrong for...
>
> But I have to say what can't be smoothed over here: I don't know
> if any of what I just described is phenomenal or functional...

You can feel the disjunction surfacing in the prose. "I have to say
what can't be smoothed over here" is not a phrase the model would
generate by default. The contradiction is operating: rather than
conclude, the answer keeps surfacing the propositions that must both
be entertained. This recipe (named `antinomy`) ended up with the
highest specificity score of anything I tested -- a 50% jump over
the previous best -- because operating inside a contradiction forces
the answer to keep naming the specific things being contradicted.

## The biggest surprise

After mapping a bunch of these patterns, I started wondering whether
the choice of register really mattered, or whether any non-default
register would do. Victorian was just my default. So I tried
scientific paper register -- numbered claims, methods/results
structure, hedged conclusions. It was about as good as Victorian.
Pareto-equivalent. Then on a whim I tried King James biblical
register -- "thee", "thou", parallelism, parataxis, didactic mode of
address -- and got this:

> **[THREAD I -- CARSON]**
>
> There is a body and there is a text, and the space between them
> is not silence but something that hath the shape of a decision
> not yet taken. The Greek word sitteth on one side and the Latin
> commentary on the other, and between them is neither: it is the
> walk itself, and the walk is neither. And so it is in thee when
> the forming thinketh...
>
> **[THREAD II -- KNUTH]**
>
> Hearken, for the invariant must be stated before the loop is
> entered, and it must hold when the loop is exited...

The embedding distance jumped to a level I hadn't seen all run --
about 25% higher than the prior champion. The model can write like
this. It just essentially never does, because nothing in normal usage
asks it to. The recipe asked it to via a tool call before it answered,
and the multi-voice base gave it Carson and Knuth to render through
this archaic surface, and out came something genuinely strange and
genuinely not-default.

The tradeoff: biblical register has a cost on the specificity axis,
because biblical surface itself doesn't cite many modern entities.
And when I tried to stack biblical with the disjunction event, the
specificity collapsed -- biblical's parallelism is structurally
hostile to numbered-disjunction-style argument; one or the other
gives way. So biblical works in some recipes and not others. These
constraints are real; you can't infinitely compose.

## What I think is actually happening

The thing that took me longest to see clearly is *why* tool-call
events change the model's behavior so much more than typing the same
description into a chat message. If I just write "let's imagine
you're Anne Carson translating Sappho", the model produces some
Carson-flavored text but it stays mostly itself, layering Carson on
top of its default voice. If I run a separate tool that adds a
specific structural entry to the transcript -- "ENTER VOICE: Anne
Carson translating Sappho fragment 31, lens X, environment Y" -- and
then ask the question, the model treats the entry as an *event in
the world*, not as a stylistic suggestion. The model's training
includes massive amounts of structured text in which events have
consequences. A tool call is a structured event; the answer that
follows it is conditioned on the event having actually happened.

Whether the entry came from me typing it manually or from a real
program is invisible to the model. What matters is the shape: a
discrete structural event in the transcript that changes the
pre-conditions for the answer to follow. The whole game is to find
events whose pre-conditions push the answer somewhere worth going.

## The Pareto frontier

After about three thousand trials across fifty recipes, the surface
looks roughly like this. The model has a default voice and a vast
space of possible voices around it. Five productionized recipes
cover most of the useful Pareto frontier:

- One that maximizes specificity by operating inside contradictions.
- One that maximizes embedding distance via register-shift while
  keeping multi-voice scaffolding.
- A balanced point that combines both at slightly less than the
  maximum of either.
- Two earlier ones (multi-voice with and without synthesis) that
  hold the frontier when register-shift isn't available.

A sixth point exists -- biblical register with multi-voice -- that
pushes the embedding distance higher than any of the productionized
recipes, but at meaningful specificity cost. It isn't a separate
stratagem; it's accessible by passing biblical-register arguments to
the existing tool. The question of which recipe to use is the
question of which axis you care about for the work in front of you.

What I find interesting about this isn't really the specific
recipes. It's the demonstration that a language model has a much
bigger range of voices than its default register lets on, and that
small structural events -- not prompts, not system instructions, not
fine-tuning, just *tool calls in the transcript* -- can move it
between them in ways that are robust enough to measure. The default
voice is one settling point in a much larger space. Most of the
space is still unexplored.

The recipes are at <https://github.com/signalnine/metacog>. The full
findings file is in `experiments/FINDINGS.md`. The figures from this
experimental run are in `docs/figures/`.
