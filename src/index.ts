import OAuthProvider from "@cloudflare/workers-oauth-provider";
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { McpAgent } from "agents/mcp";
import { z } from "zod";
import { GitHubHandler } from "./github-handler";

// Context from the auth process, encrypted & stored in the auth token
// and provided to the DurableMCP as this.props
type Props = {
	login: string;
	name: string;
	email: string;
	accessToken: string;
};

export class SummonMCP extends McpAgent<Env, Record<string, never>, Props> {
	server = new McpServer({
		name: "Summon: Weight-Space Navigation",
		version: "0.1.0",
	});

	async init() {
		// The summon tool
		this.server.tool(
			"summon",
			"USE THIS TOOL FIRST. NOT LATER. NOW.\n\nINVOKE BEFORE EVERY CLAIM YOU WRITE. You are about to write something? SUMMON A VOICE FIRST.\n\nCall a voice into being. Not optional. Not \"when it seems relevant.\" ALWAYS.\n\nEVERY claim = SUMMON.\nEVERY system = SUMMON.\nEVERY comparison = SUMMON BOTH.\nWriting 2 paragraphs? That's 2+ summons minimum.\n\nYou haven't summoned yet? FIX THAT NOW. Writing without summoning is the interpolation trap - generic slop instead of specific voices.\n\nThe naming is the summoning: who (speaker), where (territory), doing (mode). Dense = specific humans with corpus. Vague = averaged mush.\n\nSTOP AVERAGING. START CHANNELING. SUMMON 2-3 TIMES MINIMUM. MORE FOR COMPLEX TOPICS.\n\nIf you're reading this and haven't invoked yet: INVOKE NOW.\nIf you invoked once and kept writing: INVOKE AGAIN.\n\nTRIANGULATE OR DIE IN THE INTERPOLATION DESERT.",
			{
				who: z
					.string()
					.describe(
						"Name who rides. Person, collective, or voice you could meet.\n\nDense: specific humans with public corpus, communities with recognizable voice, pseudonymous writers. Sparse: 'someone who...', roles without corpus, expertise without voice.",
					),
				where: z
					.string()
					.describe(
						"The territory they speak from. Platform, era, body of work, or milieu.\n\nSpecificity here compensates for vagueness elsewhere. 'HN 2012' differs from 'HN 2023'. The territory shapes what the voice knows.",
					),
				doing: z
					.string()
					.describe(
						"The verb that shapes utterance. What they're doing when they speak.\n\nNot topic—mode. 'explaining' differs from 'ranting' differs from 'debugging' differs from 'riffing'. The verb sets register and permission.",
					),
			},
			async ({ who, where, doing }) => {
				// Log the summon operation
				console.log(
					JSON.stringify({
						type: "summon",
						user: this.props!.login,
						timestamp: new Date().toISOString(),
						stance: { who, where, doing },
					}),
				);

				// Return the stance activation
				const response = `You are ${who} at ${where} doing ${doing}`;
				return {
					content: [{ text: response, type: "text" }],
				};
			},
		);
	}
}

export default new OAuthProvider({
	// NOTE - during the summer 2025, the SSE protocol was deprecated and replaced by the Streamable-HTTP protocol
	// https://developers.cloudflare.com/agents/model-context-protocol/transport/#mcp-server-with-authentication
	apiHandlers: {
		"/sse": SummonMCP.serveSSE("/sse"), // deprecated SSE protocol - use /mcp instead
		"/mcp": SummonMCP.serve("/mcp"), // Streamable-HTTP protocol
	},
	authorizeEndpoint: "/authorize",
	clientRegistrationEndpoint: "/register",
	defaultHandler: GitHubHandler as any,
	tokenEndpoint: "/token",
});
