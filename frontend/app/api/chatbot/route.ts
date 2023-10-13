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
	const body: {
		message: string;
	} = await req.json();

	if (!body.message) {
		return new Response(JSON.stringify({ message: "Bad request" }), {
			status: 400,
			headers: { "Content-Type": "application/json" },
		});
	}

	let endpoint: string;
	let bodyToServer: string;

	try {
		const res = await postOpenAI([{ role: "user", content: body.message }]);

		const parsedRes = await res.json();
		endpoint = parsedRes.choices[0].message.function_call.name;
		bodyToServer = parsedRes.choices[0].message.function_call.arguments;
	} catch (error) {
		console.error("Error in postOpenAI:", error);
		return new Response(
			JSON.stringify({ message: "Internal Server Error" }),
			{
				status: 500,
				headers: { "Content-Type": "application/json" },
			},
		);
	}

	let serverResult: any;

	try {
		const res = await postToServer(endpoint, bodyToServer);

		serverResult = await res.json();
	} catch (error) {
		console.error("Error in postOpenAI:", error);
		return new Response(
			JSON.stringify({ message: "Internal Server Error" }),
			{
				status: 500,
				headers: { "Content-Type": "application/json" },
			},
		);
	}

	try {
		const res = await postOpenAI([
			{ role: "user", content: body.message },
			{
				role: "assistant",
				content: null,
				function_call: {
					name: endpoint,
					arguments: bodyToServer,
				},
			},
			{
				role: "function",
				name: endpoint,
				content: JSON.stringify(serverResult),
			},
		]);

		const parsedRes = await res.json();
		return new Response(
			JSON.stringify({ message: parsedRes.choices[0].message.content }),
			{
				status: 200,
				headers: { "Content-Type": "application/json" },
			},
		);
	} catch (error) {
		console.error("Error in postOpenAI:", error);
		return new Response(
			JSON.stringify({ message: "Internal Server Error" }),
			{
				status: 500,
				headers: { "Content-Type": "application/json" },
			},
		);
	}
}
