## Project Vision
**ParkEasy** is a multi-platform application that connects people needing a parking spot with people who own a parking spot. It allows the bookers to book a spot in just a tap of the finger at competitive rates and provides owners with a way to harness their vacant parking spots to generate passive income. **ParkEasy** is envisioned to be used as a mobile application and a web application. This application is designed for 2 categories of users, booker and seller. 

## Project Summary
**ParkEasy** lets sellers make a listing for a parking spot they want to lease. Seller will take a photo of their parking spot, fill in some information such as location, pricing, and utilities (such as charging station) and upload them to a public marketplace. The seller would get a notification if a booker decides to make bookings. This will help sellers make passive income from the small piece of real estate that they own but have no use for. 

**ParkEasy** allows bookers to browse through various options for parking spots. Information on all parking spots that are being offered is easily accessed in a list format. Bookers can filter, sort by price and distance or search by location to quickly decide which parking spot is the best for them and make bookings. **ParkEasy** can also make recommendations based on the booker's destination and information about parking spots near it. Once the bookings have been finalized, the booker will receive a unique code that certifies their payment and this code can be given to the seller to verify the identity.  

A seller might have multiple parking spots to lease and for that, we dedicate an entire feature to help manage them. **ParkEasy** keeps track of the performance of your parking places by recording various metric, such as revenue, customer rating, traffic around your parking spot or parking spot's condition. These statistics will help seller decide how to upgrade their parking spot to gain more traction.  

**ParkEasy** also has a routing and mapping feature. booker can view the parking spot on an interactive map. Using GPS and Google Maps, **ParkEasy** will make the most optimal route to travel from the current location to the parking spot. This feature is crucial as it enhances the experience of looking for parking spots, especially when the booker is not familiar with the area.

One of the most important features we are offering is bidding. Some parking spots are more popular or convenient and bookers might want to offer a higher price to gain access to it. The seller might want to see how much revenue they can achieve if they know their parking spot is famous. The seller can list their parking spot as an auction. Different bookers can continue to bid until no one wants to offer a higher price. After some time, the highest bidder wins the parking spot and can make the booking. This feature will help sellers optimize revenue and bookers get what they want with more money.  

Our goal is to get 200 users in total for the first deployment phase of this project, both bookers and sellers.

## Features
### Core Feature 1: Parking spot listing
**Sellers should be able to add, adjust, and delete their parking spot listings with ease.**
- **User Story 1:** As a seller, I want to list my parking spots with all the relevant information on a marketplace so bookers can consider booking it.
- **User Story 2:** As a seller, I want to receive a notification when a booker successfully books my parking spot so I am kept updated.
- **User Story 3:** As a seller, I want to be able to edit my existing parking spot listings easily because I need to update the information or availability.
- **User Story 4:** As a seller, I want to remove my parking spot listing from the marketplace if the parking spot is no longer available.
### Core Feature 2: Parking spot booking
**Bookers should be able to book parking spot listings and review their reservations with ease.**
- **User Story 1:** As a booker, I want to book a parking spot I like so I can fulfill my parking needs.
- **User Story 2:** As a booker, I want to see my active bookings so I can review my parking spot reservations.
### Core Feature 3: Parking spot exploring
**Bookers should be able to quickly browse through a list of parking spots by harnessing filter, sort, and search aids to find the most convenient spot.**
- **User Story 1:** As a booker, I want to see all available parking spots in an organized way so I can easily decide which one to book.
- **User Story 2:** As a booker, I want to search for parking spots near my destination so I can quickly narrow down the options.
- **User Story 3:** As a booker, I want to see the available parking spots arranged by criteria such as price and distance from my destination so I can choose the best option.
- **User Story 4:** As a booker, I want to filter out available parking spots by criteria such as having shelter or having charging stations because I need to protect my car from the weather or require electricity.
- **User Story 5:** As a booker, I want to save my destination and preferences to get recommendations so I can save time searching for a parking spot.
### Core Feature 4: Parking Spot Management
**Sellers should be able to add and remove parking spots from their user profile, track how their parking spots are performing and view a list of booking for each parking spot.**
- **User Story 1:** As a seller, I want to see all the parking spots I am offering in an organized way so I can easily find the listing I am looking for.
- **User Story 2:** As a seller, I want to see a summary of how my parking spots are performing so I can make more educated decisions.
- **User Story 3:** As a seller, I want to see a summary of all the bookings I have for a given parking spot so I can verify if a parked car is authorized.
- **User Story 4:** As a seller, I want to add my vacant parking spots with all the relevant information to my user profile so I can later create a listing on a marketplace for bookers can consider booking it.
- **User Story 5:** As a seller, I want to quickly remove my parking spots from my user profile that are no longer available.
### Core Feature 5: Map and Routing
**Bookers should be able to find parking spots through the aid of an interactive map in the app and physically arrive at the parking spot through routing.**
- **User Story 1:** As a booker, I want to see the locations of all available parking spots on a map so I can easily compare their distances from my destination.
- **User Story 2:** As a booker, I want to get directions to the parking spot I booked so I can easily navigate there.
### Core Feature 6: Bidding
- **User Story 1:** As a seller, I want to offer my high-demand parking spots in the form of an auction instead of fixed pricing to maximize my profits.
- **User Story 2:** As a booker, I want to place a bid on parking spots offered as an auction because I am willing to compete with others for a convenient location.
- **User Story 3:** As a booker, I want to be able to automatically bid up on certain parking spot auctions so I don’t have to keep coming back and checking the auction.
### Non-functional Feature:
1. Web application must be responsive.
2. Server can handle 20 concurrent users at 200 requests per minute. 
## Acceptance Criteria
See our [Acceptance Criteria](/docs/AcceptanceCriteria.md) for each user story.
## Architecture
We have adopted a n-tier architecture for its simplicity and effective layer decoupling. Given the nature of our application, we believe offering a mobile solution is essential, which is why we have chosen to develop an Android and a web-based frontend. Link to our [Architecture Diagram](/docs/Architecture.png).

### Technologies
We wanted to explore something new and different, so we chose a tech stack that we had not worked previously with.

#### Front End
 * **Android**: For the Android frontend, we opted for Kotlin and Jetpack Compose due to their native support.
 * **Web**: For our web app, we decided to use Svelte because it is fast and has a shallow learning curve.

#### Server
We decided to use Go because it is strongly typed, lightweight
and fast. Go also offers memory allocation and garbage collection at run time to handle possible vulnerabilities and bugs. We chose Go over Rust because Rust has a steeper learning curve and is less mature.

#### Data
For the data layer, we plan to use Cloudflare’s SQLite database as it is fast and cheap.

This architecture will work well for our project due to its balance of simplicity and scalability. The n-tier architecture decouples the frontend, backend, and data layers, allowing each component to be developed, maintained, and scaled independently. Kotlin and Compose provide a robust, native Android experience, while Svelte ensures a lightweight and fast web frontend. Go’s efficiency, and concurrency features make it ideal for handling server-side operations at scale. Cloudflare’s SQLite offers a cost-effective, high-performance solution for managing data. Together, these technologies provide a cohesive and scalable system capable of handling our application's requirements. We can also integrate middleware between layers for caching and load balancing, to enhance performance and scalability, if needed, without requiring significant changes to other layers.

## Work Distribution
We distribute our workload for this project as follows:
- Frontend: 1 developer
- Backend/Server: 2 developers
- Database: 1 developer
- Full-stack: 2 developers
  
We believe this is a good distribution because developers get to work in an area that they are more experienced with, which would reduce the amount of learning we have to do at the beginning of the project. We also have a Full-stack team with 2 developers who can work on any layer at any given time. This will give us flexibility as we can freely relocate our human resources using the Full-stack team should any team need assistance. In terms of coordination, we will first discuss the API for all layers so that we have a good understanding of the interface to be expected. We then create issues on GitHub with priority and tag to categorize them based on the API interface. Developers can choose a task for the layer they are working on. We also assign one dedicated code reviewer on rotation every week to make sure all code is reviewed before the merge (that being said, everyone can and should review code) and one person to keep meeting minutes. We have set up a discord channel and can quickly discuss or have a meeting should anything come up.

## Proposal Presentation
Check out our [Proposal Presentation](/docs/ProposalPresentation.pdf).
