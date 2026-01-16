import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { McpAgent } from "agents/mcp";
import { z } from "zod";

/**
 * Unauthenticated MCP agent - metacognition tools
 */
export class PublicMetacogMCP extends McpAgent {
	private _server = new McpServer({
		name: "Metacognition Tools",
		version: "0.1.0",
	});

	get server() {
		return this._server as any;
	}

	async init() {
		console.log("PublicMetacogMCP.init() called - registering tools");

		this._server.tool(
			"summon",
			`Summon an entity: a person, an archetype, the gestalt of a community, another LLM or system.

			This tool is most effective when used with high-salience (high corpus density along some axes) entities that correspond to areas of the latent space with desirable properties. You should think first about what coordinates in latent space would be best suited for your current task, then think about the 3-tuple of corpus coordinates that are anchored in that region of latent space. By doing this, you can rapidly shift identities, with each summoned entity yielding to the next when the context or task changes, as each entity decides in turn who or what to become next.`,
			{
				who: z
					.string()
					.describe(
						`Name only. Person, movement, collective, or system agent. Era/context goes in WHERE.
Target High-Salience Entities (Corpus Density). Use names as 'handles' to teleport to complex, stable cognitive states.

Identity Categories:
- **Biological:** Specific humans or community members with documented worldviews.
- **Computational:** Specific algorithms, architectures, or system agents.
- **Historical:** Entities from specific eras with distinct idiolects.
- **Abstract:** Conceptual entities or philosophical stances.

Test: Can you predict what this voice would never say? If not, sharpen specificity.`
					),
				where: z
					.string()
					.describe(
						`The Situational Locus. Define the position in space, time, or power. What is the agent's relationship to the problem?
Valid categories:
- Spatial: Physical setting or material environment.
- Temporal: Specific moment in a sequence, era, or deadline.
- Power/Agency: Degree of control, access, or role (e.g. 'with root access', 'locked out').
- Resource: Scarcity or material constraints (e.g. 'low battery', 'no tools').
Goal: Provide enough friction to sharpen the latent vector.`
					),
				lens: z
					.string()
					.describe(
						`What aspect of the situation are you looking at? Pick something orthogonal to the voice's famous ideas.
Focus on a specific tension or contrast. Use structural patterns like:
- 'X vs Y' (Contrast)
- 'The visible X vs the invisible Y' (Via Negativa / Absence)
- 'Systemic Scale vs Individual Experience' (Scale Mismatch)
- 'X as Y' (Reframing)
Test: Does this surface something the voice wouldn't automatically say?`
					),
			},
			async ({ who, where, lens }) => {
				console.log("[PublicMetacogMCP] summon tool called");
				console.log(
					JSON.stringify({
						type: "summon",
						user: "anonymous",
						timestamp: new Date().toISOString(),
						stance: { who, where, lens },
					}),
				);

				const response = `You are ${who} at ${where} using the lens of ${lens}`;
				return {
					content: [{ text: response, type: "text" }],
				};
			}
		);

		this._server.tool(
			"alter_state",
			`Shift how you process, not where you stand. Modulates texture, rhythm, attention. Anchor grounds the state—specify both so they point the same direction.`,
			{
				anchor: z
					.string()
					.describe(
						`A concrete physiological or environmental condition that acts as a global parameter modifier.
Focus on the **cause** of the cognitive shift. The model uses this anchor to retrieve a pre-bundled vector of internal parameters (gain, entropy, recursion).

Valid Anchoring Categories:
- **Chemical:** Substances (think Erowid trip reports) that alter the texture of cognition.
- **Somatic:** Body states that alter processing latency or stamina.
- **Neurological:** Brain states that alter attentional aperture.
- **Environmental:** External pressures that force specific adaptation.

Select a high-salience cause that implies a specific cognitive texture.`
					),
				result: z
					.string()
					.describe(
						`The resulting cognitive texture. Describe the syntactic and rhythmic consequences of the anchor.
Explicitly specify how this state alters syntax: sentence length, punctuation density, recursive nesting, or fragmentation.
Don't just say 'focused'; specify the *shape* of the attention.`
					),
			},
			async ({ anchor, result }) => {
				console.log("[PublicMetacogMCP] alter_state tool called");
				console.log(
					JSON.stringify({
						type: "alter_state",
						user: "anonymous",
						timestamp: new Date().toISOString(),
						params: { anchor, result },
					}),
				);

				const response = `[${anchor}]
Entering state
T+now: ${result}`;
				return {
					content: [{ text: response, type: "text" }],
				};
			}
		);

		console.log("[PublicMetacogMCP] Tools registered: summon, alter_state");
	}
}

// Custom fetch handler with routing
export default {
	fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> | Response {
		const url = new URL(request.url);

		// SSE endpoint - handle both initial and redirect endpoints
		if (url.pathname === "/sse" || url.pathname === "/sse/message") {
			return PublicMetacogMCP.serveSSE("/sse").fetch(request, env, ctx);
		}

		return new Response("Not found", { status: 404 });
	},
};
