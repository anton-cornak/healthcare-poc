# Healthcare POC

## How to run it locally:
- Make sure you have docker and docker-compose installed
- `cd` into the root directory of the project
- add your own OpenAI API key (from [this link](https://platform.openai.com/api-keys)) to the `docker-compose.yaml` file as `OPENAI_API_KEY` without quotes
- `docker-compose up --build`
