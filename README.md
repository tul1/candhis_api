# Candhis Wave Data Scraper & API

This project is a Golang service that scrapes wave data from the [Candhis website](https://candhis.cerema.fr/), specifically from the table containing buoy data, and provides this information via a REST API. This service is built to acquire and store data of waves from multiple buoys for further usage and analysis.

## Project Overview

Candhis (Centre d'Archivage National de Données de Houle In-Situ) provides public access to wave data from buoys around the coast of France. However, no public API is currently available for direct data access, so this project aims to fill that gap by:

- Scraping the wave data table available on the [Candhis Campaigns page](https://candhis.cerema.fr/_public_/campagne.php).
- Storing the data in a local or cloud database for easy retrieval and persistence.
- Exposing the data via a simple REST API.

### Example Data

The URL [Candhis buoy data for Les Pierres Noires](https://candhis.cerema.fr/_public_/campagne.php?Y2FtcD0wMjkxMQ==) is an example of the data provided by Candhis, showing wave data collected by the buoy located at Les Pierres Noires. The scraper will extract this data and store it in a structured format for further use.

## Features

- **Data Scraper:** A script that extracts wave data (including buoys, timestamps, wave heights, etc.) from the Candhis web page and stores it in a structured format.
- **Data Storage:** The data can be stored in a persistent storage solution such as PostgreSQL, MySQL, or any preferred database.
- **REST API:** A set of API endpoints that allow external clients to query the wave data. The API provides filtered access to data based on buoy ID, time range, and other relevant parameters.
- **Automation:** The scraper can be scheduled to run periodically to fetch and update data as needed.