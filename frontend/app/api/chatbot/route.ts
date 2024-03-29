import functionDescription from "./functions";

async function postOpenAI(messages: object[]): Promise<Response> {
	const apiKey = process.env.OPENAI_API_KEY;
	const openAIURL = process.env.OPENAI_URL;

	const body = {
		model: "gpt-3.5-turbo-0613",
		messages: messages,
		functions: functionDescription,
	};

	const headers = {
		"Content-Type": "application/json",
		Authorization: `Bearer ${apiKey}`,
	};

	return fetch(`${openAIURL}`, {
		method: "POST",
		headers,
		body: JSON.stringify(body),
	});
}

async function postToServer(endpoint: string, body: string): Promise<Response> {
	const serverURL = process.env.SERVER_URL;

	const headers = {
		"Content-Type": "application/json",
	};

	return fetch(`${serverURL}/${endpoint}`, {
		method: "POST",
		headers,
		body: body,
	});
}

export async function POST(req: Request): Promise<Response> {
	const body = await req.json();

	if (!body.message) {
		return new Response(JSON.stringify({ message: "Bad request" }), {
			status: 400,
			headers: { "Content-Type": "application/json" },
		});
	}

	try {
		const openAIResponse = await postOpenAI([
			{ role: "user", content: body.message },
		]);
		const { choices } = await openAIResponse.json();
		const choice = choices[0];
		console.log;

		if (!choice || !choice.message) {
			throw new Error("Invalid response from OpenAI");
		}

		// Handle direct message responses
		if (!choice.message.function_call) {
			return new Response(
				JSON.stringify({ message: choice.message.content }),
				{
					status: 200,
					headers: { "Content-Type": "application/json" },
				},
			);
		}

		const { name: endpoint, arguments: bodyToServer } =
			choice.message.function_call;

		// Send request to the determined server endpoint
		const serverResponse = await postToServer(endpoint, bodyToServer);
		const serverResult = await serverResponse.json();

		// Return processed result back to the user
		const finalResponse = await postOpenAI([
			{ role: "user", content: body.message },
			{
				role: "assistant",
				content: null,
				function_call: { name: endpoint, arguments: bodyToServer },
			},
			{
				role: "function",
				name: endpoint,
				content: JSON.stringify(serverResult),
			},
		]);

		const finalResult = await finalResponse.json();
		return new Response(
			JSON.stringify({ message: finalResult.choices[0].message.content }),
			{
				status: 200,
				headers: { "Content-Type": "application/json" },
			},
		);
	} catch (error) {
		console.error("Error:", error);
		return new Response(
			JSON.stringify({ message: "Internal Server Error" }),
			{
				status: 500,
				headers: { "Content-Type": "application/json" },
			},
		);
	}
}
