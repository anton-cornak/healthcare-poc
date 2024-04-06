"use client";
import { useEffect, useRef, useState } from "react";
import ReactMarkdown from "react-markdown";

interface Message {
	id: number;
	text: string;
	sender: "user" | "bot";
	isMarkdown?: boolean;
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
				id: Date.now(),
				text: "Dobrý deň. Moje meno je Zidan Sufurki a som váš asistent. Mojou úlohou je pomôcť Vám pri hľadaní lekára vo vašom okolí. Na začiatok mi, prosím, napíšte, akého lekára hľadáte a kde sa nachádzate.",
				sender: "bot",
			},
		]);
		// const text =
		// 	"Našiel som pre vás niekoľko ortopédov v Košiciach: 1. **Dr. Anton Cornak** Address: Dr. Anton Cornak Phone: 707 Bone Blvd, Orthopedics City, OC 09876 Email: 654-987-0123 Website: [julia.carter@example.com](mailto:julia.carter@example.com) 2. **Dr. Antonin Cornak** Address: Dr. Antonin Cornak Phone: 707 Bone Blvd, Orthopedics City, OC 09876 Email: 654-987-0123 Website: [julia.carter@example.com](mailto:julia.carter@example.com) 3. **Dr. Matka Kresna** Address: Dr. Matka Kresna Phone: 707 Bone Blvd, Orthopedics City, OC 09876 Email: 654-987-0123 Website: [julia.carter@example.com](mailto:julia.carter@example.com) 4. **Dr. Zidan Sufurki** Address: Dr. Zidan Sufurki Phone: 707 Bone Blvd, Orthopedics City, OC 09876 Email: 654-987-0123 Website: [julia.carter@example.com](mailto:julia.carter@example.com) 5. **Dr. David Sufuski** Address: Dr. David Sufuski Phone: 707 Bone Blvd, Orthopedics City, OC 09876 Email: 654-987-0123 Website: [julia.carter@example.com](mailto:julia.carter@example.com) Máte záujem o ďalšie informácie alebo mám pre vás vykonať ďalšiu službu?";
		// setMessages([
		// 	{
		// 		id: Date.now(),
		// 		text: formatDoctorInfo(text),
		// 		sender: "bot",
		// 		isMarkdown: containsMarkdown(text),
		// 	},
		// ]);
	}, []);

	const addMessage = (
		text: string,
		sender: "user" | "bot",
		isMarkdown: boolean = false,
	) => {
		setMessages((prevMessages) => [
			...prevMessages,
			{ id: Date.now(), text, sender, isMarkdown },
		]);
	};

	const handleUserMessage = async () => {
		if (userMessage.trim() !== "") {
			addMessage(userMessage, "user");
			setUserMessage("");

			const botResponse = await fetch("/api/chatbot", {
				method: "POST",
				body: JSON.stringify({ message: userMessage }),
			}).then((response) => response.json());

			const isMarkdown = containsMarkdown(botResponse.message);
			addMessage(
				formatDoctorInfo(botResponse.message),
				"bot",
				isMarkdown,
			);
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
							key={message.id}
							className={`mb-2 flex ${
								message.sender === "user"
									? "justify-end"
									: "justify-start"
							}`}
						>
							<div
								className={`rounded px-4 py-2 ${
									message.sender === "user"
										? "bg-blue-500 text-white"
										: "bg-gray-300 text-black"
								}`}
							>
								{message.isMarkdown ? (
									<ReactMarkdown>
										{message.text}
									</ReactMarkdown>
								) : (
									message.text
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
