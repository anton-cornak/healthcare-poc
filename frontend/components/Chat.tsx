"use client";
import { useState, useRef, useEffect } from "react";

export default function Chat() {
	const [messages, setMessages] = useState<string[]>([]);
	const [userMessage, setUserMessage] = useState<string>("");
	const chatContainerRef = useRef<HTMLDivElement | null>(null);

	const handleUserMessage = async () => {
		if (userMessage.trim() !== "") {
			setMessages((prevMessages) => [
				...prevMessages,
				`You: ${userMessage}`,
			]);

			const botResponse = await fetch("/api/chatbot", {
				method: "POST",
				body: JSON.stringify({ message: userMessage }),
			}).then((response) => response.json());

			setMessages((prevMessages) => [
				...prevMessages,
				`Bot: ${botResponse.message}`,
			]);

			setUserMessage("");
		}
	};

	const handleDeleteMessage = (indexToDelete: number) => {
		const updatedMessages = messages.filter(
			(_, index) => index !== indexToDelete,
		);
		setMessages(updatedMessages);
	};

	const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
		if (e.key === "Enter") {
			handleUserMessage();
			setUserMessage("");
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
			<div className="bg-gray-200 p-4 w-[70vh]">
				<div
					className="bg-white rounded shadow-lg p-4 h-[60vh] overflow-y-auto"
					ref={chatContainerRef}
				>
					{messages.map((message, index) => (
						<div
							key={index}
							className="mb-2 flex items-center justify-between"
						>
							<div>{message}</div>
							<button
								className="text-red-600 hover:text-red-800"
								onClick={() => handleDeleteMessage(index)}
							>
								Delete
							</button>
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
						className="rounded-r-lg p-2 bg-blue-500 text-white hover:bg-blue-600"
						onClick={handleUserMessage}
					>
						Send
					</button>
				</div>
			</div>
		</div>
	);
}
