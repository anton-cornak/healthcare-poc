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
- https://www.geoportalksk.sk/mviewer/?lang=sk&config=apps/zdravotnictvo/zdravotnictvo.xml#
- https://www.geoportalksk.sk/geonetwork/srv/eng/catalog.search#/search?facet.q=inspireThemeURI%2Fhttp%253A%252F%252Finspire.ec.europa.eu%252Ftheme%252Fhh
- https://www.geoportalksk.sk/geonetwork/doc/api/index.html#/
- **SPECIALISTS URL** https://www.geoportalksk.sk/geoserver/wfs?request=GetFeature&service=WFS&version=1.1.0&typeName=ksk_evucsk:specializovane_ambulancie_ksk&outputFormat=application%2Fjson
- add google maps link to navigate to the specialist location
- add open hours of specialists, AI needs the current time to be able to recommend the specialist in case of emergency
- search for the specialist by the name
