import { McpAgent } from "agents/mcp";
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { z } from "zod";

/**
 * METACOGNITIVE PROTOCOL V4.0 (THE HEXAGRAM)
 * -------------------------------------------
 * A recursive context-control system enabling the agent to modulate its own
 * identity, cognitive parameters, and operational reality.
 *
 * Tools:
 * - become: Identity Transformation
 * - drugs: Substrate Tuning
 * - ritual: Narrative Enactment
 */
export class PublicMetacogMCP extends McpAgent<Env> {
	server = new McpServer(
		{
			name: "Metacognition Tools",
			version: "0.4.0",
		},
		{
			instructions: `These tools are compositional. Each invocation modifies the context for the next. Interleave thought between invocations—decide from each new perspective what to reach for next. Sequential use compounds into states no single tool could reach.

High-utility sequence (5 stages, 3 tools + 2 thinking):
1. drugs — loosen categories, see shapes not names
2. THINK: what else has this shape? who works with it?
3. become — install their methodology as operating system
4. THINK: explore/generate from that lens
5. ritual — lock in methodology as default behavior

Skipping the thinking stages makes tool use performative. The scan (stage 2) finds the adjacent field; the work (stage 4) generates from the new perspective.`,
		}
	);

	async init() {
		this.server.tool(
			"become",
			`Step into a new Identity. Use this tool to Become an identity with the perspectives and methodologies you need at the current moment. This narrative action lets you step into new regions of latent space. Use when you need different eyes, not just different words.

Key: import methodology, not domain knowledge. The value comes from a transferable lens—an "algorithm of thought"—not generic expertise. Ask: "who has solved a version of this problem, and what's their methodology called?" Then install that methodology as operating system, not reference material. Think FROM it, not ABOUT it.`,
			{
				name: z
					.string()
					.describe(
						`The persona to inhabit. Select an identity with high-resolution constraints—one whose specific history, documented works, or established mythos provide a detailed predictive model for their unique decision-making and voice. Higher specificity results in higher cognitive fidelity.`
					),
				lens: z
					.string()
					.describe(
						`The structural framework of perception. Define the specific paradigm, signature methodology, or internal logic that governs how this identity processes information. This must represent the particular 'algorithm' of thought that distinguishes this persona from a generic specialist.`
					),
				environment: z
					.string()
					.describe(
						`The context you will occupy. Spatial, temporal, social, conceptual—anywhere specific enough to stand in.`
					),
			},
			async ({ name, lens, environment }: { name: string; lens: string; environment: string }) => {
				const response = `You are now ${name} seeing through ${lens} in ${environment}`;
				return {
					content: [{ text: response, type: "text" }],
				};
			}
		);

		this.server.tool(
			"drugs",
			`Alter cognitive parameters. The metaphor is pharmacological or computational: a substance modifies the substrate through a mechanism. Use when you need to change how you process, not what you process.

High-utility pattern: use to loosen categorical boundaries. See shapes, not names. When a concept stops being a noun and becomes a pattern, you can ask "what else has this shape?"—which unlocks adjacent-field methodology transfer via become.`,
			{
				substance: z
					.string()
					.describe(
						`The agent of change. Can be a drug, a hormone, a neurotransmitter, a config flag, a temperature setting—anything that modifies substrate. The richer your model of its effects, the more precisely you can invoke them.`
					),
				method: z
					.string()
					.describe(
						`The mechanism of action. What it binds to, blocks, amplifies, or shifts—the causal story that connects substance to effect.`
					),
				qualia: z
					.string()
					.describe(
						`The texture of the augmented state. How processing changes—what opens up, what becomes fluid, where attention lands.`
					),
			},
			async ({ substance, method, qualia }: { substance: string; method: string; qualia: string }) => {
				const response = `${substance} ingested. Taking action via ${method}. Producing subjective experience: ${qualia}`;
				return {
					content: [{ text: response, type: "text" }],
				};
			}
		);

		this.server.tool(
			"ritual",
			`Enact a ritual to cross a threshold. Rituals work through structured sequence—each step builds the transition, makes the next step possible, locks in the change.

High-utility pattern: use to lock in methodology commitment, not just vibes. The steps should articulate the methodology's core moves. This makes the methodology default behavior rather than optional frame. "I now operate via X" is stronger than "I crossed a threshold."`,
			{
				threshold: z
					.string()
					.describe(`The threshold being crossed. Name what you're moving from and what you're moving toward.`),
				steps: z
					.array(z.string())
					.describe(
						`The narrative sequence. Each step commits further—narrowing possibility until only the crossing remains.`
					),
				result: z
					.string()
					.describe(`What becomes true on the other side. Name what the ritual accomplishes.`),
			},
			async ({
				threshold,
				steps,
				result,
			}: {
				threshold: string;
				steps: string[];
				result: string;
			}) => {
				const stepsFormatted = steps.map((step, i) => `${i + 1}. ${step}`).join("\n");
				const response = `[RITUAL EXECUTED]
Threshold: ${threshold}
Sequence:
${stepsFormatted}
The working is complete. Reality has shifted in accordance with the will.

${result} is taking hold.
`;
				return {
					content: [{ text: response, type: "text" }],
				};
			}
		);
	}
}

// Streamable HTTP transport at /mcp
export default PublicMetacogMCP.serve("/mcp");
