"use client";
import { useEffect, useRef, useState } from "react";
import ReactMarkdown from "react-markdown";

interface Message {
	content: string;
	role: "user" | "assistant";
}

function formatDoctorInfo(text: string): string {
	// Split the text by the pattern that indicates a new doctor entry, which in this case is a digit followed by a period and a space.
	const doctorEntries = text
		.split(/\d\. \*\*/)
		.filter((entry) => entry.trim() !== "");

	// Prefix that will be added before each doctor's name to maintain the markdown bold syntax after splitting
	const prefix = "**";

	// Process each entry to ensure it's formatted on a new line
	const formattedEntries = doctorEntries.map((entry, index) => {
		// Re-add the markdown bold syntax for the doctor's name
		const formattedEntry = `${index + 1}. ${prefix}${entry.trim()}`;

		// Replace the spaces after "Address:", "Phone:", "Email:", and "Website:" with newlines for readability
		return formattedEntry
			.replace(/Address:/g, "\n\n   Adresa:")
			.replace(/Phone:/g, "\n\n   Telefónne číslo:")
			.replace(/Email:/g, "\n\n   Email:")
			.replace(/Website:/g, "\n\n   Webová stránka:");
	});

	// Join the processed entries back into a single string, each entry separated by two newlines for clear separation
	return formattedEntries.join("\n\n").substring(5);
}

export default function Chat() {
	const [messages, setMessages] = useState<Message[]>([]);
	const [userMessage, setUserMessage] = useState<string>("");
	const chatContainerRef = useRef<HTMLDivElement | null>(null);

	const containsMarkdown = (text: string) =>
		text.includes("**") || text.includes("[");

	useEffect(() => {
		setMessages([
			{
				content:
					"Dobrý deň. Moje meno je Zidan Sufurki a som váš asistent. Mojou úlohou je pomôcť Vám pri hľadaní lekára vo vašom okolí. Na začiatok mi, prosím, napíšte, akého lekára hľadáte a kde sa nachádzate.",
				role: "assistant",
			},
		]);
	}, []);

	const addMessage = (content: string, role: "user" | "assistant") => {
		setMessages((prevMessages) => [...prevMessages, { content, role }]);
	};

	const handleUserMessage = async () => {
		if (userMessage.trim() !== "") {
			const newMessage: Message = {
				content: userMessage,
				role: "user",
			};
			const updatedMessages = [...messages, newMessage];
			addMessage(userMessage, "user");
			setUserMessage("");

			const botResponse = await fetch("/api/chatbot", {
				method: "POST",
				body: JSON.stringify({
					message: userMessage,
					conversation: updatedMessages,
				}),
			}).then((response) => response.json());

			addMessage(formatDoctorInfo(botResponse.message), "assistant");
		}
	};

	const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
		if (e.key === "Enter") {
			handleUserMessage();
		}
	};

	useEffect(() => {
		if (chatContainerRef.current) {
			chatContainerRef.current.scrollTop =
				chatContainerRef.current.scrollHeight;
		}
	}, [messages]);

	return (
		<div className="min-h-screen flex justify-center items-center">
			<div className="bg-gray-200 p-4 w-full max-w-md">
				<div
					className="bg-white rounded shadow-lg p-4 h-[60vh] overflow-y-auto"
					ref={chatContainerRef}
				>
					{messages.map((message) => (
						<div
							key={message.content}
							className={`mb-2 flex ${
								message.role === "user"
									? "justify-end"
									: "justify-start"
							}`}
						>
							<div
								className={`rounded px-4 py-2 ${
									message.role === "user"
										? "bg-blue-500 text-white"
										: "bg-gray-300 text-black"
								}`}
							>
								{containsMarkdown(message.content) ? (
									<ReactMarkdown>
										{message.content}
									</ReactMarkdown>
								) : (
									message.content
								)}
							</div>
						</div>
					))}
				</div>
				<div className="flex mt-4">
					<input
						type="text"
						className="flex-1 rounded-l-lg p-2 border-t mr-0 border-b border-l text-gray-800 border-gray-200 bg-white"
						placeholder="Type a message..."
						value={userMessage}
						onChange={(e) => setUserMessage(e.target.value)}
						onKeyDown={handleKeyPress}
					/>
					<button
						className="rounded-r-lg px-4 bg-blue-500 text-white hover:bg-blue-600"
						onClick={handleUserMessage}
					>
						Send
					</button>
				</div>
			</div>
		</div>
	);
}
