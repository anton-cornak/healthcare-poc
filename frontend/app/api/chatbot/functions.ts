const functionDescription = [
	{
		name: "add",
		description: "Add two or more numbers together",
		parameters: {
			type: "object",
			properties: {
				numbers: {
					type: "array",
					items: {
						type: "number",
					},
					description: "An array of numbers to be added together",
				},
			},
			required: ["numbers"],
			description: "Payload containing numbers to add",
		},
	},
	{
		name: "subtract",
		description: "Subtract two or more numbers from the base number",
		parameters: {
			type: "object",
			properties: {
				number: {
					type: "number",
					description: "The base number to subtract from",
				},
				subtract: {
					type: "array",
					items: {
						type: "number",
					},
					description:
						"An array of numbers to be substracted from the base number",
				},
			},
			required: ["number, subtract"],
			description: "Payload containing numbers to subtract from a number",
		},
	},
	{
		name: "compute",
		description: "Add or substract list of numbers",
		parameters: {
			type: "object",
			properties: {
				add: {
					type: "array",
					items: {
						type: "number",
					},
					description: "An array of numbers to be added together",
				},
				subtract: {
					type: "array",
					items: {
						type: "number",
					},
					description:
						"An array of numbers to be substracted from the added numbers or from each other if no numbers are added",
				},
			},
			required: ["number, subtract"],
			description: "Payload containing numbers to subtract from a number",
		},
	},
];

export default functionDescription;
