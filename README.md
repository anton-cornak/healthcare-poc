# Healthcare POC

## How to run it locally:
- Make sure you have docker and docker-compose installed
- `cd` into the root directory of the project
- add your own OpenAI API key (from [this link](https://platform.openai.com/api-keys)) to the `docker-compose.yaml` file as `OPENAI_API_KEY` without quotes
- add your own Geocode API KEY (from [this link](https://geocode.maps.co/join/)) to the `docker-compose.yaml` file as `GEOCODE_API_KEY` without quotes
- `docker-compose up --build`

If any changes to the database are made (especially in `init.sql`), run `docker-compose down -v` and then `docker-compose up --build` again.

### Comments:
- https://www.topdoktor.sk/hodnotenie-lekarov/
