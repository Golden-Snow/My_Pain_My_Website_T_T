#######################################################################################################################
                                     Country Info & Population Service
#######################################################################################################################


Hi there! :3
This project is a simple web service built using the Go programming language.
It helps you look up country information and historical population data by talking to a couple of external services.
The service shows details like a country’s name, population, languages, borders, and even a list of its cities! x3 
It also lets you see population numbers over time and calculates the average.


#######################################################################################################################
                                        What This Project Does
#######################################################################################################################

Country Information:
When you enter a 2-letter country code (for example, "no" for Norway or "us" for the United States), 
the service gives you a bunch of details about that country. You’ll see the country’s name,
the continent it’s in, population count, languages spoken, borders with other countries, the URL for its flag, 
and even a list of cities (which you can limit to a specific number).

Population Data:
You can also check out historical population numbers for the country. If you want, 
you can filter the results by providing a start and end year. 
The service will show you the numbers for each year and the average population for that period.

Diagnostics:
There’s a special section that checks whether the external services (that provide the country and population data) 
are working properly. It also tells you how long this service has been running.


#######################################################################################################################
                                     How the Project Is Organized
#######################################################################################################################


Navn(Folder)
|
|---> Main.go
|---> Readme
|
|---> FrontEnd(folder)
    |---> index.html
    |---> styles.css
    |---> script.js
    |
    |--->Pictures(Folder)
        |
        |--->Snowie.png

main.go:
This is the main file where all the “behind the scenes” work happens.
It handles the communication with external APIs, processes the data, and connects everything together.

FrontEnd Folder:
Here you’ll find the files for the web interface. :D

index.html: A simple web page where you can test the service.
script.js: The JavaScript that makes the API calls and displays the results on the page.
styles.css: Non existant! but you can add stuff here if you want.

Pictures Folder:
This folder contains a personal picture that I included in the project. no stealing! >:3


#######################################################################################################################
                                How to Run This Project on Your Computer
#######################################################################################################################


What do i need?

To run this project, you need to have Go installed on your computer. You can download it from the official Go website.

Link: https://go.dev/dl/

If the go version casues trouble for you then try using the Go version used when creating this:
go version: 1.24.0

----------------------------------------

I got Go installed, what now?

You can either clone the repository or download the ZIP file and extract it. they are located under these two links.

GitHub: https://github.com/Golden-Snow/My_Pain_My_Website_T_T
GitLab: https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2025-workspace/ommarkus/assignment-1_my-absolute-struggles

----------------------------------------

I got the files, and I got Go installed, now what?

Open Your Command Prompt or Terminal 
(This has many names if you dont use english on your computer, but typing "CMD" is a good way to find it) 

Now you need to navigate to the folder where the main.go file is located inside the Command Prompt. 
For example, if you extracted the folder on your Desktop, use this command:


If you downloaded from GitHub: cd Desktop/My_Pain_My_Website_T_T
If you downloaded from GitLab: cd Desktop/Assignment-1_My-absolute-struggles-main

If the name is different just put the name of the folder in when using the cd command.

----------------------------------------

I am in the folder now inside the terminal but its only a black screen? what now?

Now its time to finally start up the server!

With Go installed, type the following command:

type this in the Command Prompt: go run main.go

This will start the server on your computer! 
it usually listening to port 8080 but if you are special person it might be on a different port.
in the Command Prompt window it will say where it listening.

----------------------------------------

Ok it said its listening to a port! ok.. what do i do now?

Open your prefered web browser! the website should be up for you right now!

type this in the top websearch bar: http://localhost:8080

if it was listening to different port then put that number in instead of the 8080 in the link above.

----------------------------------------

Now enjoy!!! :3

You’ll see the simple web interface where you can enter a 2-letter country code and get the information you need.


#######################################################################################################################
                                          How to Use the Service
#######################################################################################################################


Looking Up Country Info:
Type a 2-letter country code (like "no", "us", or "nl") into the field. 
You can also set a limit on the number of cities displayed if you want. 
Then, click on "Get Country Info" and the details will appear below in a neat JSON format.

Getting Population Data:
Enter a country code along with a start and an end year if you want to filter the population data. 
Click "Get Population Data" to see the historical numbers and the average.

Checking Diagnostics:
Click on "Get Diagnostics Status" to see if the external services are working and to check the uptime of the service.


#######################################################################################################################
                                            It's Online too!
#######################################################################################################################


You can find this service online aswell!
The service is deployed on Render. 

You can try it out online here: https://my-pain-my-website-t-t.onrender.com


#######################################################################################################################
                                    Wanne Set Up An Online one too?
#######################################################################################################################


Deployment on Render!

Make an account on render: https://render.com/

once you made an account,
Log in to Render and select New > Web Service.
Connect to your GitHub repository containing this code.

Configure the Service:
Language: Go
Root Directory: Leave blank (assuming main.go is at the repo root, if not you need to make a path).

Build Command: go build -o main main.go

Start Command: ./main

Set Environment Variables: Not needed, leave blank.

Render automatically sets PORT for you :3

Click Create Web Service (or Deploy Web Service).
Then monitor the build logs until the service is successfully deployed.

Render provides a URL like https://my-country-info.onrender.com.
Visit that URL in your browser or use a REST client (e.g., Postman) to test endpoints such as /countryinfo/v1/info/no.

That’s it! Your service should now be live on Render. If you make changes to your repo, you can redeploy by enabling 
Auto-Deploy or clicking Manual Deploy in your service settings.


#######################################################################################################################
                                               API Endpoints
#######################################################################################################################


Country Info
Path: /countryinfo/v1/info/{code}?limit={cityLimit}
Method: GET
Description: Returns country details (name, continent, population, languages, borders, flag URL, capital, and cities).

Extra Note:
{code} is the mandatory 2-letter country code.
limit (optional) restricts the number of cities (sorted alphabetically).
Implementation: Uses fetchRestCountry() for official data and fetchCities() for the cities list.

-----------------------------------------------------------------------------------------------------------------------

Population Data
Path: /countryinfo/v1/population/{code}?limit=YYYY-YYYY
Method: GET
Description: Provides historical population counts for a country and calculates the average.

Extra Note:
{code} is the mandatory 2-letter country code.
limit (optional) specifies a year range (e.g., 2010-2015).
Implementation: Uses fetchRestCountry() for the official name and fetchPopulation() for historical data, 
then filters and computes the mean.

-----------------------------------------------------------------------------------------------------------------------

Diagnostics
Path: /countryinfo/v1/status/
Method: GET
Description: Shows the HTTP status of external APIs, the service version, and uptime.
Implementation: Uses getServiceStatus() to check external API health and reports the uptime.


#######################################################################################################################
                                           External Dependencies
#######################################################################################################################


REST Countries API

Endpoint: http://129.241.150.113:8080/v3.1/
Usage: Retrieves official country data (name, population, capital, etc.).
Code: Called via fetchRestCountry().


CountriesNow API

Endpoint: http://129.241.150.113:3500/api/v0.1/
Usage: Provides city lists and historical population data.
Code: Accessed by fetchCities() for cities and fetchPopulation() for population counts.


#######################################################################################################################
                                               License:
#######################################################################################################################


All code in this project is free to use, modify, and distribute. However, personal pictures included in the project 
are the exclusive property of the author and may not be used without permission. :3


#######################################################################################################################
                                              Final Note!
#######################################################################################################################


This project took many long nights of coding and troubleshooting to get it right. 
I used simple HTML for the frontend to keep things easy and clear. Although the CSS is very basic, 
the main goal was to ensure the backend (main.go) could effectively process and display the data you need. 
There are safeguards in place to handle unusual inputs (like incorrect year ranges) so that the service doesn’t 
break easily. ^w^

If you wanne contribute or have any questions or just wanne provide some feedback! 
feel free to reach out to me on discord:

New tag: @goldensnow
Old tag: Golden Snow#7489

