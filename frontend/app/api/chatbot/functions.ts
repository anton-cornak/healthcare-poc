// OpenAI returns 'specialties/all' does not match '^[a-zA-Z0-9_-]{1,64}$'
const functionMapping: { [key: string]: string } = {
	"location-wkt": "location/wkt",
	"specialties-all": "specialties/all",
	"specialist-find": "specialist/find",
};

const functionDescription = [
	{
		name: "location-wkt",
		description:
			"Generates WKT representation of a location from a string representation of the location",
		parameters: {
			type: "object",
			properties: {
				user_location: {
					type: "string",
					description:
						"A string representation of the user's location, e.g. 'London, UK'",
				},
			},
			required: ["user_location"],
			description:
				"Payload containing user location as a string representation",
		},
	},
	{
		name: "specialties-all",
		description: "Gets list of all specialties in the database",
	},
	{
		name: "specialist-find",
		description:
			"Gets a specialist from the database by the user defined preferences, such as specialty, location and distance from the user",
		parameters: {
			type: "object",
			properties: {
				specialty_id: {
					type: "number",
					description:
						"Name of the specialty from the list of specialties - this should always correspond to the /specialties endpoint response",
				},
				radius: {
					type: "number",
					description:
						"Radius a user is willing to travel to see a specialist. Radius should always be in METERS.",
				},
				user_location: {
					type: "string",
					description:
						"WKT representation representation of the user's location, e.g. 'POINT(21.2496774 48.7172272)'",
				},
			},
			required: ["specialty_id", "radius", "user_location"],
			description:
				"Payload containing user preferences when searching for a specialist",
		},
	},
];

export { functionDescription, functionMapping };
