"use client";
import { useEffect, useRef, useState } from "react";

interface Message {
	id: number;
	text: string;
	sender: "user" | "bot";
}

export default function Chat() {
	const [messages, setMessages] = useState<Message[]>([]);
	const [userMessage, setUserMessage] = useState<string>("");
	const chatContainerRef = useRef<HTMLDivElement | null>(null);

	useEffect(() => {
		setMessages([
			{
				id: Date.now(),
				text: "Dobrý deň. Moje meno je Zidan Sufurki a som váš asistent. Mojou úlohou je pomôcť Vám pri hľadaní lekára vo vašom okolí. Na začiatok mi, prosím, napíšte, akého lekára hľadáte a kde sa nachádzate.",
				sender: "bot",
			},
		]);
	}, []);

	const addMessage = (text: string, sender: "user" | "bot") => {
		setMessages((prevMessages) => [
			...prevMessages,
			{ id: Date.now(), text, sender },
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

			addMessage(botResponse.message, "bot");
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
								{message.text}
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
