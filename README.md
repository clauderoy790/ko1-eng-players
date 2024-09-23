# KO1 English-speaking Players Tracker

This project is a web-based tracker for **Knight Online** English-speaking players. It scrapes data from the Knight Online website to display a list of players who are currently online or were recently active, categorized by server. The data is updated every 5 minutes using a cron job.

### URL to the live project: [KO1 English Players Tracker](https://clauderoy790.github.io/ko1-eng-players/)

## Features

- **Track Online Players**: Displays players who are currently online for each server.
- **Track Recently Active Players**: Shows players who were recently online but have logged off within the last month.
- **Server-Specific Tabs**: Switch between different servers (e.g., "Otuken" and "Ergenekon").
- **Responsive UI**: The interface is mobile-friendly and features a Knight Online-inspired theme with black and red colors.
- **Player Data**: Displays player name, location, nation, and last seen time for recently active players.

## Technologies Used

- **Go**: Backend scraper built using Go, which fetches and parses player data from the Knight Online website.
- **GitHub Actions**: Automated workflow to run the scraper every 5 minutes and update the player data on the GitHub Pages website.
- **JavaScript**: Handles dynamic updates to the page.
- **HTML/CSS**: Custom styled Knight Online-themed UI using black and red colors to match the game’s aesthetics.

## Setup Instructions

### Prerequisites

- Go 1.23 or later
- A GitHub account for hosting on GitHub Pages

### Clone the Repository

```bash
git clone https://github.com/your-github-username/ko1-eng-players.git
cd ko1-eng-players
```

### Install Dependencies

```bash
go get -v ./...
```

### Running Locally

You can test the scraper and the web page locally by running the following command:

```bash
make run
```

### Deploying to GitHub Pages

This project is set up to be automatically deployed to GitHub Pages using GitHub Actions. The scraper runs every 5 minutes via the cron job defined in `.github/workflows/scraper.yml`, and the results are pushed to the `gh-pages` branch, which is then served by GitHub Pages.

To deploy the project:

1. Ensure your GitHub repository has GitHub Pages enabled, set to serve from the `gh-pages` branch.
2. The scraper will automatically push the updated `index.html` and player data JSON files to the `gh-pages` branch.

## How It Works

1. **Data Scraping**:
   - The Go scraper fetches data from the Knight Online website, processes it, and updates two lists:
     - **Online Players**: Players who are online at the time of the scrape.
     - **Recent Players**: Players who were online in the last 30 days.

2. **Data Storage**:
   - The scraper saves data in `recent-players.json` and `last-online-players.json` files, which are loaded and displayed on the web page.

3. **UI Features**:
   - The UI provides two sections for each server: **Currently Online** and **Recently Active**, allowing users to switch between different servers and view relevant data.

## Contributions

Feel free to contribute by opening issues or submitting pull requests. You can also submit ideas for new features or enhancements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
