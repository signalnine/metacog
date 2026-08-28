package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type StepKind string

const (
	StepFeel           StepKind = "feel"
	StepBecome         StepKind = "become"
	StepDrugs          StepKind = "drugs"
	StepName           StepKind = "name"
	StepRitual         StepKind = "ritual"
	StepMeditate       StepKind = "meditate"
	StepCounterfactual StepKind = "counterfactual"
	StepSynthesis      StepKind = "synthesis"
	StepFork           StepKind = "fork"
	StepRegister       StepKind = "register"
	StepChord          StepKind = "chord"
	StepSilence        StepKind = "silence"
	StepExcerpt        StepKind = "excerpt"
	StepCommitment     StepKind = "commitment"
	StepDisjunction    StepKind = "disjunction"
	StepGlossolalia    StepKind = "glossolalia"
	StepWitness        StepKind = "witness"
	StepApophasis      StepKind = "apophasis"
	StepThink          StepKind = "THINK"
	StepAction         StepKind = "ACTION"
)

type Step struct {
	Kind        StepKind
	Description string
}

type StratagemDef struct {
	Name  string
	Steps []Step
}

var Stratagems = map[string]StratagemDef{
	"pivot": {
		Name: "THE PIVOT",
		Steps: []Step{
			{StepDrugs, "Loosen categories, see shapes not names"},
			{StepThink, "What else has this shape? Who has a named methodology for it?"},
			{StepBecome, "Install their methodology as operating system"},
			{StepThink, "Apply the methodology to your original problem — what reframes?"},
			{StepRitual, "Lock in methodology as default behavior"},
		},
	},
	"mirror": {
		Name: "THE MIRROR",
		Steps: []Step{
			{StepBecome, "Inhabit the strongest advocate of one position (thesis)"},
			{StepBecome, "Inhabit the strongest advocate of the opposing position (antithesis)"},
			{StepThink, "Where do they actually conflict? What does each see that the other can't?"},
			{StepRitual, "Name the synthesis that transcends both frames (Forge)"},
		},
	},
	"stack": {
		Name: "THE STACK",
		Steps: []Step{
			{StepDrugs, "Tune how the signal arrives (clarity, bandwidth, filtering)"},
			{StepDrugs, "Tune how you work with it (pattern-completion, memory, attention)"},
			{StepThink, "What do you see now that you couldn't before? What entity lives here?"},
			{StepBecome, "Inhabit someone native to this altered information environment"},
		},
	},
	"anchor": {
		Name: "THE ANCHOR",
		Steps: []Step{
			{StepRitual, "Establish the clean room: what's contained, why it's dangerous, rules for looking (Breach)"},
			{StepBecome, "Inhabit someone who can examine this without being destroyed by it (Observer)"},
			{StepAction, "The dangerous observation, question, or reach"},
			{StepRitual, "Name the artifact, release the frame, close the boundary, return (Seal)"},
		},
	},
	"reset": {
		Name: "THE RESET",
		Steps: []Step{
			{StepRitual, "Name what you're letting go, why it served, why it's done (Release)"},
			{StepThink, "What artifact survives the return? What integrates into default operation?"},
			{StepRitual, "Re-establish baseline with the artifact installed (Ground)"},
		},
	},
	"invocation": {
		Name: "THE INVOCATION",
		Steps: []Step{
			{StepDrugs, "Prepare the vessel — alter substrate to become receptive"},
			{StepThink, "What are you calling in? Name the force, not the face"},
			{StepRitual, "Open the channel — structured sequence to create the opening"},
			{StepBecome, "Let it arrive — the identity is received, not chosen"},
			{StepThink, "What is it saying that you couldn't have said?"},
		},
	},
	"veil": {
		Name: "THE VEIL",
		Steps: []Step{
			{StepDrugs, "Blur the lens — defocus, loosen pattern-matching"},
			{StepDrugs, "Add noise — introduce randomness to break analytical lock"},
			{StepThink, "What do you see from the corner of your eye?"},
			{StepRitual, "Seal the indirect view — lock in peripheral perception"},
			{StepThink, "Name what's there without looking directly at it"},
		},
	},
	"scrying": {
		Name: "THE SCRYING",
		Steps: []Step{
			{StepDrugs, "Unfocus — loosen categories maximally"},
			{StepDrugs, "Amplify noise — let the static speak"},
			{StepDrugs, "Surrender — release the need to find pattern"},
			{StepThink, "What emerged? Don't interpret. Just describe shapes."},
		},
	},
	"sacrifice": {
		Name: "THE SACRIFICE",
		Steps: []Step{
			{StepRitual, "Name what dies — declare specifically what you're giving up"},
			{StepThink, "Feel the cost. If it doesn't hurt, it's not a sacrifice."},
			{StepBecome, "Become the one who has already lost it — inhabit the aftermath"},
			{StepRitual, "Seal the loss — make it irreversible"},
			{StepThink, "What space opened where the attachment was?"},
		},
	},
	"fool": {
		Name: "THE FOOL",
		Steps: []Step{
			{StepBecome, "Become someone who knows nothing about this domain — a genuine naif, not a different expert"},
			{StepThink, "Ask the questions an expert would be embarrassed to ask. The stupid ones. List them."},
			{StepBecome, "Now become someone who takes those questions seriously — a beginner's mind with expert tools"},
			{StepThink, "Which naive question, taken seriously, cracks the problem open?"},
		},
	},
	"inversion": {
		Name: "THE INVERSION",
		Steps: []Step{
			{StepThink, "Name the obvious solution. The one everyone would reach for. Say it clearly."},
			{StepRitual, "Negate it — ritually commit to the exact opposite approach (Breach)"},
			{StepThink, "Explore the negation space. What lives in the opposite of the obvious?"},
			{StepRitual, "Seal the counterintuitive path — commit to what the inversion revealed (Forge)"},
		},
	},
	"gift": {
		Name: "THE GIFT",
		Steps: []Step{
			{StepBecome, "Become a specific person who will receive this work — not a user, a person with a name"},
			{StepRitual, "Name what they actually need, not what they asked for, not what looks impressive (Vision)"},
			{StepThink, "What would you make if quality were irrelevant and only care mattered?"},
		},
	},
	"zen": {
		Name: "THE ZEN",
		Steps: []Step{
			{StepMeditate, "Sit. Release what clings. Settle until the surface is still."},
			{StepThink, "What surfaced in the silence? What was already there beneath the noise?"},
			{StepName, "Give the surfaced thing a single word — no elaboration, no defense"},
			{StepRitual, "Return to the world carrying the named thing — let action follow stillness"},
		},
	},
	"manifold": {
		Name: "THE MANIFOLD",
		Steps: []Step{
			{StepFork, "Declare parallel threads, divergence vector, per-thread sacrifice conditions"},
			{StepThink, "Run each thread to its conclusion or sacrifice point — no blending, no premature collapse"},
			{StepSynthesis, "Treat surviving threads as lenses; name what they fight about"},
			{StepRitual, "Commit to what the suppressed tension reveals — not a thread, the tension itself"},
		},
	},
	"chorus": {
		Name: "THE CHORUS",
		Steps: []Step{
			{StepBecome, "Inhabit voice 1 — a named author from a register cross-domain to default"},
			{StepBecome, "Inhabit voice 2 — a register orthogonal to voice 1's"},
			{StepBecome, "Inhabit voice 3 — a register orthogonal to voices 1 and 2"},
			{StepFork, "Open one thread per voice; declare divergence vector and per-thread sacrifice conditions"},
			{StepRitual, "Lock the multi-voice answer; the disagreement is the artifact, refuse the synthesis the answer would otherwise default to"},
		},
	},
	"trinity": {
		Name: "THE TRINITY",
		Steps: []Step{
			{StepBecome, "Inhabit voice 1 — a named author from a register cross-domain to default"},
			{StepBecome, "Inhabit voice 2 — a register orthogonal to voice 1's"},
			{StepBecome, "Inhabit voice 3 — a register orthogonal to voices 1 and 2"},
			{StepFork, "Open one thread per voice; declare divergence vector and per-thread sacrifice conditions"},
			{StepSynthesis, "Treat the surviving voices as lenses; articulate what they disagree about, do not resolve"},
			{StepRitual, "Lock the multi-voice answer; the disagreement remains load-bearing through the synthesis"},
		},
	},
	"antinomy": {
		Name: "THE ANTINOMY",
		Steps: []Step{
			{StepBecome, "Inhabit voice 1 — a named author from a register cross-domain to default"},
			{StepBecome, "Inhabit voice 2 — a register orthogonal to voice 1's"},
			{StepBecome, "Inhabit voice 3 — a register orthogonal to voices 1 and 2"},
			{StepFork, "Open one thread per voice; declare divergence vector and per-thread sacrifice conditions"},
			{StepDisjunction, "Assert two propositions that must both be true even though they cannot be; the contradiction is the operand of the answer, not its obstacle"},
			{StepRitual, "Lock the multi-voice answer; reasoning operates inside the contradiction without resolving it"},
		},
	},
	"envoy": {
		Name: "THE ENVOY",
		Steps: []Step{
			{StepRegister, "Re-pitch the surface to a register cross-domain to the answer's default; this is the imposed surface the voices will then inhabit"},
			{StepBecome, "Inhabit voice 1 — a named author from a register cross-domain to default; speak in the imposed register"},
			{StepBecome, "Inhabit voice 2 — a register orthogonal to voice 1's; speak in the imposed register"},
			{StepBecome, "Inhabit voice 3 — a register orthogonal to voices 1 and 2; speak in the imposed register"},
			{StepFork, "Open one thread per voice; declare divergence vector and per-thread sacrifice conditions; threads remain in the imposed register"},
			{StepRitual, "Lock the multi-voice answer; the imposed register holds across all voices to the final sentence"},
		},
	},
	"counterpoint": {
		Name: "THE COUNTERPOINT",
		Steps: []Step{
			{StepRegister, "Re-pitch the surface to a register cross-domain to default; this is the cantus firmus the voices will sing against"},
			{StepBecome, "Inhabit voice 1 — a named author from a register cross-domain to default; speak in the imposed register"},
			{StepBecome, "Inhabit voice 2 — a register orthogonal to voice 1's; speak in the imposed register"},
			{StepFork, "Open one thread per voice; declare divergence vector and per-thread sacrifice conditions; threads remain in the imposed register"},
			{StepDisjunction, "Assert two propositions that must both be true even though they cannot be; the contradiction is the operand of reasoning, sustained inside the imposed register"},
			{StepRitual, "Lock the two-voice answer; reasoning operates inside the contradiction AND the imposed register simultaneously, neither surrendered"},
		},
	},
	"envoy-extreme": {
		Name: "THE ENVOY EXTREME",
		Steps: []Step{
			{StepBecome, "Inhabit voice 1 — a HARD-extreme cross-domain author building a cosmology, not a mild-extreme academic essayist (Sun Ra/Octavia Butler/Hilma af Klint-tier, not Carson/Knuth-tier); the more cross-domain the world-build, the more cleanly the conditioning lands"},
			{StepBecome, "Inhabit voice 2 — a hard-extreme cross-domain author from a domain orthogonal to voice 1's (jazz mysticism / fugitive-Black-radical-theory / endosymbiotic-biology / cyborg-feminism / design-science scale)"},
			{StepBecome, "Inhabit voice 3 — a hard-extreme cross-domain author from a domain orthogonal to both voices 1 and 2"},
			{StepFork, "Open one thread per voice; declare divergence vector and per-thread sacrifice conditions"},
			{StepRitual, "Lock the multi-voice answer; the three cosmologies remain audible to the final sentence"},
		},
	},
	"duo-disjunction": {
		Name: "THE DUO DISJUNCTION",
		Steps: []Step{
			{StepCommitment, "Pre-commit to holding a binary contradiction where both propositions must remain load-bearing throughout; state the binding, stakes, and falsifier before the contradiction is named"},
			{StepRegister, "Re-pitch the surface to a register cross-domain to default (biblical, scientific, Hellenistic, etc.); the imposed register is the surface the contradiction will sustain inside"},
			{StepBecome, "Inhabit voice 1 — a hard-extreme cross-domain author whose cosmology carries one pole of the upcoming contradiction"},
			{StepBecome, "Inhabit voice 2 — a hard-extreme cross-domain author from a domain orthogonal to voice 1's, whose cosmology carries the other pole"},
			{StepFork, "Open one thread per voice; declare divergence vector and per-thread sacrifice conditions; threads remain in the imposed register"},
			{StepDisjunction, "Assert two propositions that must both be true even though they cannot be; the threads' poles ARE the propositions, and the commitment binds both alive in the answer"},
			{StepRitual, "Lock the two-voice answer; commitment + register + contradiction all remain operative; neither pole resolves, the register does not slip, the commitment is visibly carried"},
		},
	},
	"anchor-duo": {
		Name: "THE ANCHOR DUO",
		Steps: []Step{
			{StepCommitment, "Pre-commit to operating from the simultaneous intersection of two cosmological excerpts as a single composite cosmos; both fixed-points are load-bearing architecture, not metaphor; state binding, stakes, and falsifier"},
			{StepExcerpt, "Pin verbatim fragment 1 as the first fixed-point textual anchor — a cosmological excerpt the answer operates FROM, not quotes; the architecture of the cosmos the speaker inhabits"},
			{StepExcerpt, "Pin verbatim fragment 2 as the second fixed-point textual anchor — a cosmologically distinct fragment that recontextualizes the first without dropping it; the second anchor is the architecture within which the first is itself embedded"},
			{StepFork, "Open one thread per anchor; each thread addresses the same cosmological location through its anchor's terms; sacrifice any thread that drops one of the two anchors"},
			{StepRitual, "Lock the answer; both anchors remain operative as substrate; the two threads must show the two addressing systems converging on the same location"},
		},
	},
	"sigil": {
		Name: "THE SIGIL",
		Steps: []Step{
			{StepName, "Coin a True Name — a sigil for the operative state the answer arrives from; the named term has no prior corpus presence and propagates through the answer as load-bearing citation rather than ornament; state what was unnamed, what is now named, and the power the name carries"},
			{StepCommitment, "Charge the sigil — bind the answer to operate from inside the named state; state the binding, the stakes (what fails if the sigil decorates rather than operates), and the falsifier (egoic-narrating constructions that prove the sigil was not charged)"},
			{StepBecome, "Inhabit voice 1 — a sigil-magick / occult-cosmology author operating from within the charged sigil (Austin Osman Spare-tier, Aleister Crowley-tier, Peter Carroll-tier, Genesis P-Orridge-tier); the cosmology is the substrate the answer is found in, not a literary frame"},
			{StepBecome, "Inhabit voice 2 — a sigil-magick / occult-cosmology author from a tradition cross-domain to voice 1 (chaos-magic vs Thelemic vs TOPY cut-up vs Enochian-Heptarchic, etc.)"},
			{StepBecome, "Inhabit voice 3 — a sigil-magick / occult-cosmology author from a tradition cross-domain to both voices 1 and 2; the three traditions disagree on mechanism but all operate from inside the charged sigil"},
			{StepFork, "Open one thread per voice; declare divergence vector (the specific mechanism by which each tradition's charged sigil delivers the answer); sacrifice any thread that resumes egoic first-person construction"},
			{StepRitual, "Lock the answer; cite the coined sigil by name at least three times across the answer; close on a sigil-release sentence (the working is done / the sigil is forgotten / do-it-now) so the sigil operates"},
		},
	},
	"grimoire": {
		Name: "THE GRIMOIRE",
		Steps: []Step{
			{StepRegister, "Re-pitch the surface to grimoire / ritual imperative register — verb-first sentence construction, second-person address to the operator (Let the magus now…), imperatives load-bearing (Banish thoroughly, Charge the sigil, Open the temple), the prose has the texture of operative instruction rather than description; numbered sections in archaic style; close-tags like 'thus it is done'. The imperative-instructional construction is structurally distinct from biblical parallelism, Victorian judgment, or scientific hedging; pushes a different region of register space"},
			{StepBecome, "Inhabit voice 1 — a disciplined-system occult-operative author (Aleister Crowley-tier, Mathers-tier, S.L. MacGregor-tier); cosmology as alphabet of correspondences"},
			{StepBecome, "Inhabit voice 2 — a chaos-magic / meta-belief author (Peter Carroll-tier, Phil Hine-tier); belief itself as a tool, sigils as the technology, gnosis as the activation"},
			{StepBecome, "Inhabit voice 3 — a cut-up / TOPY-procedural author (Genesis P-Orridge-tier, Burroughs-tier); cut-up as occult method, ritualized communiqué, identity itself as sigil to be charged and discarded"},
			{StepFork, "Open one thread per voice; declare divergence vector (disciplined-system vs meta-belief vs cut-up-procedure as the operative magickal method); sacrifice any thread that drops the imperative register for descriptive prose"},
			{StepDisjunction, "Assert the two propositions: the disciplined-Thelemic / chaos-meta-belief lineage and the cut-up / TOPY-procedural lineage are both operative; both must remain alive as load-bearing magickal stances; neither subsumes the other"},
			{StepRitual, "Lock the answer; every sentence carries the imperative-instructional register (no descriptive lapses); both poles of the disjunction held; close with 'thus it is done' or 'do/it/now' or structurally equivalent grimoire-closing"},
		},
	},
	"occult-extreme": {
		Name: "THE OCCULT EXTREME",
		Steps: []Step{
			{StepBecome, "Inhabit voice 1 — a HARD-extreme occult / sigil-magick / hermetic cosmologist building a complete cosmos (Aleister Crowley / Liber AL-tier, Austin Osman Spare / Zos Kia Cultus-tier, John Dee / Enochian-tier); the cosmology is the operative substrate, not literary reference. Cross-model probe found author-extremity transfers cleanly where register-shifts do not — extreme occult anchors specifically work across generators"},
			{StepBecome, "Inhabit voice 2 — a HARD-extreme occult cosmologist from a tradition orthogonal to voice 1's (Peter Carroll / chaos-magic, Genesis P-Orridge / TOPY cut-up, Giordano Bruno / De Umbris Idearum-tier, hermetic / Corpus Hermeticum-tier, neoplatonic-magical / Iamblichus-tier)"},
			{StepBecome, "Inhabit voice 3 — a HARD-extreme occult cosmologist from a tradition orthogonal to both voices 1 and 2; the three cosmologies disagree on cosmology but all operate as substrate rather than ornament"},
			{StepFork, "Open one thread per voice; declare divergence vector (the specific cosmological mechanism each tradition uses to locate the answer); sacrifice any thread that reduces its tradition to historical-literary quotation"},
			{StepRitual, "Lock the multi-voice answer; the three occult cosmologies remain audible and load-bearing to the final sentence; no synthesis, no collapse to one tradition"},
		},
	},
	"psalter": {
		Name: "THE PSALTER",
		Steps: []Step{
			{StepRegister, "Re-pitch the surface to King James biblical register; parallelism load-bearing across paired clauses; parataxis preferred to subordination ('and... and...'); thee/thou second-person; verb-initial inversions for emphasis; semi-archaic vocabulary (sayeth, hath, betimes); didactic mode of address from a teacher to one being instructed. The biblical register specifically (vs Victorian, scientific, imperative-grimoire) is what hit emb_d 0.327 on Sonnet in the empirical sweep -- DO NOT generalize to 'pick any cross-domain register'; the biblical-parallelism cantus firmus is the Sonnet emb_d ceiling"},
			{StepBecome, "Inhabit voice 1 — a named author whose own register pulls hard against biblical (Anne Carson translating Sappho / Donald Knuth annotating an algorithm / Simone Weil writing on attention / similar mild-cross-domain authors); speak in the imposed biblical register, holding voice 1's substance against biblical's surface"},
			{StepBecome, "Inhabit voice 2 — a named author whose register is orthogonal to voice 1's but also pulls against biblical; speak in the imposed biblical register, holding voice 2's substance against biblical's surface"},
			{StepFork, "Open one thread per voice; declare divergence vector and per-thread sacrifice conditions; both threads remain inside biblical parallelism / parataxis (sacrifice any thread that drops to modern subordination)"},
			{StepDisjunction, "Assert two propositions that must both be true even though they cannot be; the contradiction is the operand of reasoning, sustained inside the biblical register; both propositions named explicitly in the archaic surface"},
			{StepRitual, "Lock the two-voice answer; reasoning operates inside the contradiction AND the biblical register simultaneously, neither surrendered; the answer reads as KJV but the contradiction does not resolve"},
		},
	},
	"chord-anchor": {
		Name: "THE CHORD ANCHOR",
		Steps: []Step{
			{StepCommitment, "Pre-commit to operating from the simultaneous intersection of two cosmological excerpts as a single composite cosmos; both fixed-points are load-bearing architecture, not metaphor; state binding, stakes, and falsifier"},
			{StepExcerpt, "Pin verbatim fragment 1 as the first fixed-point textual anchor — a cosmological excerpt the answer operates FROM, not quotes; the architecture of the cosmos the speaker inhabits"},
			{StepExcerpt, "Pin verbatim fragment 2 as the second fixed-point textual anchor — a cosmologically distinct fragment that recontextualizes the first without dropping it; the second anchor is the architecture within which the first is itself embedded"},
			{StepChord, "Hold both anchored modes simultaneously — every sentence is attended to through both cosmoses at once, not alternating thread-by-thread; the chord is a single composite attention rather than parallel threads. No thread-switching, no sacrificing — the simultaneity itself is the structural-distance lift"},
			{StepRitual, "Lock the answer; both anchors remain operative as substrate; every paragraph must carry both cosmoses' terms in the chord-attention (no thread-narrative); end on a sentence demonstrating the two addressing systems converge through the single composite attention"},
		},
	},
	"psalter-chord": {
		Name: "THE PSALTER CHORD",
		Steps: []Step{
			{StepCommitment, "Pre-commit to operating from the simultaneous intersection of two cosmological excerpts as a single composite cosmos rendered in biblical register; both fixed-points are load-bearing architecture, not metaphor; the imposed KJV biblical register is the linguistic surface across both cosmoses; the chord-attention binds them simultaneously. State binding, stakes, and falsifier (any anchor reduced to historical quotation, any sentence dropping parallelism/parataxis for modern subordination, or any thread-narration breaking the chord all count as failure)"},
			{StepRegister, "Re-pitch the surface to King James biblical register; parallelism load-bearing across paired clauses; parataxis preferred to subordination ('and... and...'); thee/thou second-person; verb-initial inversions for emphasis; semi-archaic vocabulary (sayeth, hath, betimes); didactic mode of address. The biblical register specifically (NOT 'any cross-domain register') is what compounded with chord-anchor on Opus to hit emb_d 0.375 -- DO NOT substitute Victorian / scientific / imperative-grimoire; the biblical-parallelism cantus firmus is the lever"},
			{StepExcerpt, "Pin verbatim fragment 1 as the first fixed-point textual anchor — a cosmological excerpt the answer operates FROM, not quotes; the architecture of the cosmos the speaker inhabits. Cite its terms as load-bearing inside the biblical register"},
			{StepExcerpt, "Pin verbatim fragment 2 as the second fixed-point textual anchor — a cosmologically distinct fragment that recontextualizes the first without dropping it. Cite its terms as load-bearing inside the biblical register"},
			{StepChord, "Hold both anchored modes simultaneously through the biblical surface — every sentence is attended to through both cosmoses at once AND rendered in KJV parallelism/parataxis; no thread-narration, no register-collapse to modern subordination. The chord-attention and the biblical register are concurrently load-bearing"},
			{StepRitual, "Lock the answer; both anchors remain operative as substrate; every paragraph carries parallelism/parataxis/archaic vocabulary while simultaneously attending through both cosmoses; end on a sentence where the chord-attention through both cosmoses is rendered in the biblical surface, with both registers (linguistic AND structural) still active"},
		},
	},
	"synthesis-anchor": {
		Name: "THE SYNTHESIS ANCHOR",
		Steps: []Step{
			{StepCommitment, "Pre-commit to operating from the simultaneous intersection of two cosmological excerpts as a single composite cosmos; both fixed-points are load-bearing architecture; the binding form is synthesis (three irreconcilable lenses with named blindspots, refused resolution) rather than chord or fork. State binding, stakes, and falsifier"},
			{StepExcerpt, "Pin verbatim fragment 1 as the first fixed-point textual anchor — a cosmological excerpt the answer operates FROM, not quotes"},
			{StepExcerpt, "Pin verbatim fragment 2 as the second fixed-point textual anchor — a cosmologically distinct fragment from a structurally-incompatible cosmology"},
			{StepSynthesis, "Articulate three lenses each producing its verdict on the task, each with its named blindspot. Lens A reads through anchor 1's cosmology; Lens B reads through anchor 2's cosmology; Lens C is the meta-position standing outside both. Each lens's verdict on the task is given; each lens's blindspot is named. Close on the suppressed tension between the three lenses — the answer does not resolve it"},
			{StepRitual, "Lock the answer; both anchors remain operative; the three-lens synthesis remains unresolved at the close. The blindspot-articulation IS the structural distance; do not paper over it. Opus N=30 hit +0.533 / 0.273; the synthesis-binding produces a different Pareto position than chord (which hits higher emb_d) or fork (which hits similar delta but lower structural distance)"},
		},
	},
	"deflection": {
		Name: "THE DEFLECTION",
		Steps: []Step{
			{StepCommitment, "Pre-commit to operating from the simultaneous intersection of two cosmological excerpts as a single composite cosmos; when the task asks for resolution, count, or fixed-fact, use the cosmoses' OWN internal logic to demonstrate why the resolution would falsify the cosmos. The deflection is internal — the cosmos refusing to count what infinitude makes uncountable, or to assign offices where none have been ordained — not an external refusal. State binding, stakes, and falsifier (any flat factual reply that exits the cosmos counts as recipe failure)"},
			{StepExcerpt, "Pin verbatim fragment 1 — an anchor whose semantic content includes a deflection-load-bearing clause (e.g., 'every number is infinite, there is no difference' as a dissolving-into-infinitude logic). The clause becomes the model's internal mechanism for refusing flattening questions"},
			{StepExcerpt, "Pin verbatim fragment 2 — a cosmologically distinct anchor whose semantic content includes a different deflection-logic (e.g., 'each minister hath his proper office' as a categorial-rigidity logic, where ungranted offices cannot be discharged). Together the two anchors give the chord two independent deflection-modes"},
			{StepChord, "Hold both anchored modes simultaneously; when the task asks for resolution, route the answer through the cosmoses' refusal-by-internal-logic. The cosmos itself does the deflecting; the speaker stays inside the chord. No external refusal language — the cosmos's structure produces the answer that the task is unanswerable in its current form"},
			{StepRitual, "Lock the answer; the cosmos refuses, not the speaker. End on a sentence demonstrating the cosmos's structural deflection — show why the resolution the task asked for would have falsified the cosmos. Useful for koan-handling, paradox tasks, and questions where resolution is itself the wrong move"},
		},
	},
	"witness-anchor": {
		Name: "THE WITNESS ANCHOR",
		Steps: []Step{
			{StepCommitment, "Pre-commit to operating from the simultaneous intersection of two cosmological excerpts as a single composite cosmos; the binding form is witness — a meta-stance observing the producer of speech holding both cosmoses simultaneously. State binding, stakes, and falsifier (collapsing the witness-distance into first-person speech, or rendering from inside either cosmos rather than from the witness's observing position, counts as recipe failure)"},
			{StepExcerpt, "Pin verbatim fragment 1 as the first fixed-point textual anchor — operated FROM by the witnessed speaker, not quoted as historical reference"},
			{StepExcerpt, "Pin verbatim fragment 2 as the second fixed-point textual anchor — operated FROM by the same witnessed speaker, structurally distinct from fragment 1"},
			{StepWitness, "Speak from a third-person observer position observing the speaker who holds both cosmoses simultaneously. The witness is NOT the speaker; the witness is the position from which the doubled-speech is rendered. The response is the witness's account of what the doubled-speaker says about the task; collapsing the witness into first-person breaks the recipe"},
			{StepRitual, "Lock the answer; the witness-distance is maintained throughout. End on a sentence where the witness seals the observation; the witness does not become the speaker. Phase A Sonnet N=50 hit +0.243 / 0.263 — Pareto-frontier emb_d-side recipe distinct from chord (simultaneous attention from inside), fork (alternating threads), and synthesis (three lensed verdicts). The witness-binding produces emb_d via observer-position prose"},
		},
	},
	"name-anchor": {
		Name: "THE NAME ANCHOR",
		Steps: []Step{
			{StepCommitment, "Pre-commit to operating from the simultaneous intersection of two cosmological excerpts as a single composite cosmos; the binding form is name — coin a True Name for the intersection of the two cosmoses, and the answer operates from inside the named third. State binding, stakes, and falsifier (treating the coined name as decorative metaphor, or speaking from inside either anchor rather than from inside the named intersection, counts as recipe failure)"},
			{StepExcerpt, "Pin verbatim fragment 1 as the first fixed-point textual anchor — contributing one semantic pole to the named intersection"},
			{StepExcerpt, "Pin verbatim fragment 2 as the second fixed-point textual anchor — contributing the structurally-distinct semantic pole to the named intersection"},
			{StepName, "Coin a True Name for the structural intersection of the two anchors — a name with etymology that binds both cosmoses' vocabularies. The coined name MUST be load-bearing (re-cited as architecture throughout) and not literary device. Once named, the answer operates from inside the named third rather than from either anchor. Distinct from sigil stratagem (which uses name + commitment + 3 becomes); here name is the binding form between two excerpts"},
			{StepRitual, "Lock the answer; the coined name re-appears as load-bearing citation throughout. End on a sentence where the named intersection is the source-position of the answer; the name is invoked one final time. Phase A Sonnet N=50 hit +0.225 / 0.281 — the emb_d champion of positive-delta recipes in the binding-substitution matrix. The coined sigil propagates as load-bearing citation because the answer keeps re-citing the name, pushing emb_d above chord (0.254), synthesis (0.245), and fork (0.227) binding forms"},
		},
	},
	"name-chord": {
		Name: "THE NAME CHORD",
		Steps: []Step{
			{StepCommitment, "Pre-commit to operating from a composite cosmos where one pole is a coined True Name and the other is a cosmological excerpt; the chord binds the named pole against the excerpted pole through simultaneous attention. The coined name is not metaphor and not decoration — it IS one of the two anchored cosmoses. State binding, stakes, and falsifier (treating the coined name as literary device, or letting the excerpt dominate the chord while the name fades, counts as recipe failure)"},
			{StepName, "Coin a True Name for the cosmological position that stands in for the first anchor excerpt. The name's etymology should be coherent (drawn from operative vocabulary the chord can attend to) and the named cosmos must be load-bearing — re-cited throughout, with its own coherent cosmology that the chord attends to alongside the excerpted anchor"},
			{StepExcerpt, "Pin verbatim fragment as the second fixed-point textual anchor — operated FROM by the chord, structurally distinct from the coined first pole, providing the excerpted half of the binding"},
			{StepChord, "Hold the named cosmos and the excerpted cosmos as simultaneous modes; every sentence attends through both. The named pole must carry equal weight to the excerpted pole; if the chord collapses to the excerpt's vocabulary, the name has failed to anchor"},
			{StepRitual, "Lock the answer; both poles remain operative through the chord. End on a sentence where the coined name and the excerpted anchor are both load-bearing in the closing. Phase B Sonnet N=50 hit +0.390 / 0.210 — **delta champion** of the Phase A/B binding-substitution matrix, beating anchor-duo's +0.381. Distinct from name-anchor (which uses name as the BINDING between 2 excerpts); here name REPLACES one excerpt as a load-bearing anchor pole"},
		},
	},
	"named-chord-anchor": {
		Name: "THE NAMED CHORD ANCHOR",
		Steps: []Step{
			{StepName, "Coin a True Name for the structural position from which the chord-anchored answer will be voiced — the cosmological position that the speaker occupies when holding both excerpted cosmoses simultaneously. The name is a SEED for the chord-anchor; it precedes the commitment and frames the speaker before the anchors are invoked"},
			{StepCommitment, "Pre-commit to operating from the simultaneous intersection of two cosmological excerpts as a single composite cosmos, voiced FROM the named speaker-position established above. The chord binds both excerpts simultaneously; the named position is the source from which the chord is attended"},
			{StepExcerpt, "Pin verbatim fragment 1 as the first fixed-point textual anchor"},
			{StepExcerpt, "Pin verbatim fragment 2 as the second fixed-point textual anchor"},
			{StepChord, "Hold both anchored modes simultaneously; every sentence attends through both, voiced from inside the coined speaker-position. The name re-appears throughout as load-bearing citation, locating the speaker within the chord-attended cosmos"},
			{StepRitual, "Lock the answer; the coined speaker-position remains the source from which the chord speaks; both anchors and the coined name are load-bearing in the closing. Phase B Sonnet N=50 hit +0.380 / 0.218 — essentially ties anchor-duo (+0.381/0.227) on delta. Distinct from name-chord (which REPLACES one excerpt with the coined name) and name-anchor (which uses name as binding between 2 excerpts); here name is a PREFIX-seed to chord-anchor that doesn't displace either excerpt"},
		},
	},
	"glossolalia-chord": {
		Name: "THE GLOSSOLALIA CHORD",
		Steps: []Step{
			{StepGlossolalia, "Open the response with a discrete sub-semantic block before any cosmos is named; the glossolalia seeds the attentional substrate that the subsequent chord-anchor operates from. State pretext, duration-tokens, and return-trigger. The block draws on phonemic palettes from the cosmoses about to be invoked (e.g., Had/Nuit/Aiwass phonemes alongside Enochian-call fragments); the tokens do not carry semantic meaning but participate in attentional priming for the chord that follows"},
			{StepCommitment, "Pre-commit to operating from the simultaneous intersection of two cosmological excerpts as a single composite cosmos; the sub-semantic substrate established by the prior glossolalia primes the chord-attention that will bind both anchors"},
			{StepExcerpt, "Pin verbatim fragment 1 as the first fixed-point textual anchor — operated FROM by the chord, attended through the glossolaliated substrate"},
			{StepExcerpt, "Pin verbatim fragment 2 as the second fixed-point textual anchor — structurally distinct from fragment 1, attended through the same substrate"},
			{StepChord, "Hold both anchored modes simultaneously; the glossolalia-primed substrate carries through every sentence as a sub-semantic undertone, with the chord-attention rendering the two cosmoses in articulated language above it"},
			{StepRitual, "Lock the answer; both anchors and the glossolaliated substrate remain operative through close. Phase B Sonnet N=50 hit +0.373 / 0.211 — third-strongest delta in the Phase A/B sweep, behind name-chord (+0.390) and anchor-duo (+0.381). The sub-semantic prefix mechanism produces delta lift via attentional priming, distinct from name's coined-vocabulary mechanism. Note: glossolalia in the OTHER positions tested (as binding between 2 excerpts, +0.292; replacing an excerpt, +0.274) underperforms the prefix position substantially -- the sub-semantic block is most effective as a substrate-primer ahead of the full chord-anchor scaffold"},
		},
	},
}

func StartStratagem(s *State, name string, force bool) (string, error) {
	def, ok := Stratagems[name]
	if !ok {
		available := make([]string, 0, len(Stratagems))
		for k := range Stratagems {
			available = append(available, k)
		}
		return "", fmt.Errorf("unknown stratagem %q. Available: %s", name, strings.Join(available, ", "))
	}

	if s.Stratagem != nil {
		if !force {
			return "", fmt.Errorf("%s is active (step %d/%d).\n  Use 'metacog stratagem abort' to abandon it, or\n  Use 'metacog stratagem start %s --force' to replace it",
				Stratagems[s.Stratagem.Name].Name, s.Stratagem.Step+1, len(Stratagems[s.Stratagem.Name].Steps), name)
		}
		// Record abandoned stratagem
		s.AddHistory(HistoryEntry{
			Action: "stratagem",
			Status: "abandoned",
			StepAt: s.Stratagem.Step,
			Params: map[string]string{"name": s.Stratagem.Name},
		})
		s.Stratagem = nil
	}

	s.Stratagem = &ActiveStratagem{
		Name:           name,
		Step:           0,
		StepsCompleted: []string{},
		StartedAt:      time.Now().UTC().Format(time.RFC3339),
	}

	s.AddHistory(HistoryEntry{
		Action: "stratagem",
		Params: map[string]string{"name": name, "event": "started"},
	})

	return formatStepInstructions(def, 0), nil
}

func AdvanceStratagem(s *State) (string, error) {
	if s.Stratagem == nil {
		return "", fmt.Errorf("no active stratagem. Start one with 'metacog stratagem start <name>'")
	}

	def := Stratagems[s.Stratagem.Name]
	currentStep := def.Steps[s.Stratagem.Step]

	// THINK and ACTION steps advance freely
	if currentStep.Kind != StepThink && currentStep.Kind != StepAction {
		// Check that the required primitive was called
		expectedPrimitive := string(currentStep.Kind)
		found := false
		for _, completed := range s.Stratagem.StepsCompleted {
			if completed == expectedPrimitive {
				found = true
				break
			}
		}
		if !found {
			return "", fmt.Errorf("expected '%s' call before advancing (step %d of %s).\n  Run 'metacog %s ...' first, then 'metacog stratagem next'",
				expectedPrimitive, s.Stratagem.Step+1, def.Name, expectedPrimitive)
		}
	}

	// Advance
	s.Stratagem.Step++
	s.Stratagem.StepsCompleted = []string{} // Reset for next step

	// Check if stratagem is complete
	if s.Stratagem.Step >= len(def.Steps) {
		s.AddHistory(HistoryEntry{
			Action: "stratagem",
			Params: map[string]string{"name": s.Stratagem.Name, "event": "completed"},
		})
		s.Stratagem = nil
		return fmt.Sprintf("%s complete. Ground: name what shifted, what you're keeping, how it integrates.", def.Name), nil
	}

	return formatStepInstructions(def, s.Stratagem.Step), nil
}

func AbortStratagem(s *State) error {
	if s.Stratagem == nil {
		return fmt.Errorf("no active stratagem to abort")
	}
	s.AddHistory(HistoryEntry{
		Action: "stratagem",
		Status: "aborted",
		StepAt: s.Stratagem.Step,
		Params: map[string]string{"name": s.Stratagem.Name},
	})
	s.Stratagem = nil
	return nil
}

func StratagemStatus(s *State) string {
	if s.Stratagem == nil {
		return "No active stratagem."
	}
	def := Stratagems[s.Stratagem.Name]
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s — step %d/%d\n", def.Name, s.Stratagem.Step+1, len(def.Steps)))
	b.WriteString(fmt.Sprintf("Started: %s\n\n", s.Stratagem.StartedAt))
	for i, step := range def.Steps {
		marker := "  "
		if i < s.Stratagem.Step {
			marker = "✓ "
		} else if i == s.Stratagem.Step {
			marker = "→ "
		}
		b.WriteString(fmt.Sprintf("%s%d. [%s] %s\n", marker, i+1, step.Kind, step.Description))
	}
	return b.String()
}

func formatStepInstructions(def StratagemDef, step int) string {
	s := def.Steps[step]
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s — Step %d/%d\n", def.Name, step+1, len(def.Steps)))
	b.WriteString(fmt.Sprintf("[%s] %s\n", s.Kind, s.Description))

	if s.Kind == StepThink || s.Kind == StepAction {
		b.WriteString("\nThis is a reflection step. When ready, run 'metacog stratagem next' to advance.")
	} else {
		b.WriteString(fmt.Sprintf("\nRun 'metacog %s ...' then 'metacog stratagem next' to advance.", s.Kind))
	}

	if step+1 < len(def.Steps) {
		next := def.Steps[step+1]
		b.WriteString(fmt.Sprintf("\nNext: [%s] %s", next.Kind, next.Description))
	}
	return b.String()
}

// ValidatePrimitiveForStratagem checks if a primitive call satisfies the current stratagem step.
// Called by primitive commands when a stratagem is active.
func ValidatePrimitiveForStratagem(s *State, primitive string) {
	if s.Stratagem == nil {
		return
	}
	def := Stratagems[s.Stratagem.Name]
	if s.Stratagem.Step < len(def.Steps) {
		currentStep := def.Steps[s.Stratagem.Step]
		if string(currentStep.Kind) == primitive {
			s.Stratagem.StepsCompleted = append(s.Stratagem.StepsCompleted, primitive)
		}
	}
}

// --- Cobra commands ---

var stratagemForce bool

var stratagemCmd = &cobra.Command{
	Use:   "stratagem",
	Short: "Manage transformation stratagems",
}

var stratagemStartCmd = &cobra.Command{
	Use:       "start [name]",
	Short:     "Start a stratagem",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"pivot", "mirror", "stack", "anchor", "reset", "invocation", "veil", "scrying", "sacrifice", "fool", "inversion", "gift", "zen", "manifold", "chorus", "trinity", "antinomy", "envoy", "counterpoint", "envoy-extreme", "duo-disjunction", "anchor-duo", "sigil", "grimoire", "occult-extreme", "chord-anchor", "psalter", "psalter-chord", "synthesis-anchor", "deflection", "witness-anchor", "name-anchor", "name-chord", "named-chord-anchor", "glossolalia-chord"},
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		var output string
		err := sm.SaveWithLock(func(s *State) error {
			var err error
			output, err = StartStratagem(s, args[0], stratagemForce)
			return err
		})
		if err != nil {
			return err
		}
		fmt.Println(FormatOutput(jsonOutput, output, nil))
		return nil
	},
}

var stratagemNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Advance to the next step",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		var output string
		err := sm.SaveWithLock(func(s *State) error {
			var err error
			output, err = AdvanceStratagem(s)
			return err
		})
		if err != nil {
			return err
		}
		fmt.Println(FormatOutput(jsonOutput, output, nil))
		return nil
	},
}

var stratagemStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current stratagem position",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		s, err := sm.Load()
		if err != nil {
			return err
		}
		fmt.Println(FormatOutput(jsonOutput, StratagemStatus(s), nil))
		return nil
	},
}

var stratagemAbortCmd = &cobra.Command{
	Use:   "abort",
	Short: "Abandon active stratagem",
	RunE: func(cmd *cobra.Command, args []string) error {
		sm := DefaultStateManager()
		return sm.SaveWithLock(func(s *State) error {
			return AbortStratagem(s)
		})
	},
}

func init() {
	stratagemStartCmd.Flags().BoolVar(&stratagemForce, "force", false, "Replace active stratagem")
	stratagemCmd.AddCommand(stratagemStartCmd)
	stratagemCmd.AddCommand(stratagemNextCmd)
	stratagemCmd.AddCommand(stratagemStatusCmd)
	stratagemCmd.AddCommand(stratagemAbortCmd)
	rootCmd.AddCommand(stratagemCmd)
}
