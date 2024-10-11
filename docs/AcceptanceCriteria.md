# Acceptance Criteria
### Core Feature 1: Listing
**Sellers should be able to add, adjust, and delete their parking spot listings with ease.**

- **User Story 1:** As a seller, I want to list my parking spots with all the relevant information on a marketplace so bookers can consider booking it.
	- **Scenario:** 
“Given that I am a user who has a parking spot available for booking  
**When** I open my Personal Parking Spots page, the system shows a list of my active parking spot listings with tags saying “listed” or “auctioned” for a parking spot that is listed or auctioned respectively  
**And** when I click on a parking spot without any tag, the system shows a Parking Spot Details page with the information for that parking spot, and a button that says “List” and a large text that says “Unlisted”, a listing table, and an unmodifiable “Price/hour” field that is empty  
**And** when I click on the “List” button, the system displays a listing table, “Price/Hour” field, and the "Cancel" and "Submit" buttons  
Where the listing table’s column is 7 week days from Monday to Sunday and the column is time from 7:00 am to 12:00 pm and each row is 1 hour; green cell means a listed slot, gray cell means an unlisted slot, and blue cell means a booked slot  
**And** since the parking spot I selected is not listed, the system shows all cells in the table as gray
**And** when I click on a gray cell in the listing table, the system colors the cell green to say the slot selected for listing  
**And** when I click on a green cell in the listing table, the system colors the cell gray to say the slot is unlisted   
**And** If I confirm the listing time slot that I want, fill in a valid value in the “Price/Hour” field and click the "Submit" button  
**Then** the system asks for confirmation with a small pop-up message saying “Are you sure you wish to list your parking spot?” and "Yes" and "No" buttons  
&emsp;**And** if I click “Yes”, **Then** the system submits my listing request  
&emsp;**And** the "Price/Hour" fields are updated with the new information,  
&emsp;**And** the large text that says “Unlisted” is replaced by the table with information about the time slot that I have entered  
&emsp;**And** the system displays a flash message saying “Parking spot listed successfully”  
&emsp;**And** when I open the Listing Page and search for my parking spot, the system shows my parking spot listing”  
&emsp;**Or** if I click “No”, the small pop up message will disappear and the system displays the listing table and the “Price/hour” field in the same state before I pressed the “Cancel” button  
**Or** if I click on the "Cancel" button, the system displays a message “Are you sure you wish to cancel without saving any changes?” and "Yes" and "No" buttons  
&emsp;**And** if I click "Yes", the system takes me back to the current Parking Spot Details page in the same state before I pressed the “List” button  
&emsp;**Or** if I click “No”, the small pop up message will disappear and the system displays the listing table and the Price/hour field in the same state as before I pressed the “Cancel” button”

- **User Story 2:** As a seller, I want to receive a notification when a booker successfully books my parking spot so I am kept updated.
	- **Scenario:** 
“Given that I am a user who has an active parking spot listing,  
**When** a booker books my parking spot, the system sends a notification to my phone including the booker's name and the name of the listing they booked  
**And** when I click on the notification, the system takes me to my Confirmed Bookings page of that particular parking spot which shows a list of all bookings with the booking corresponding to the notification highlighted”

- **User Story 3:** As a seller, I want to be able to edit my existing parking spot listings easily because I need to update the information or availability.
	- **Scenario:** 
“Given that I am a user who has a parking spot available for booking **When** I open my Personal Parking Spots page, the system shows a list of my active parking spot listings with tags saying “listed” or “auctioned” for a parking spot that is listed or auctioned  
**And** when I click on a listed parking spot, the system shows a Parking Spot Details page with the information for that parking spot and a button that says “Edit Listing”, a listing table, and the current price/hour  
**Where** the listing table’s column is 7 week days from Monday to Sunday and the column is time from 7:00 am to 12:00 pm and each row is 1 hour; green cell means a listed slot, gray cell means an unlisted slot, and blue cell means a booked slot  
**And** if I click the “Edit” button, the system makes the listing table and “Price/Hour” field editable and displays "Cancel" and "Update" buttons  
**When** I click on a gray cell in the listing table, the system colors the cell green to say the slot is up for listing  
**Or** when I click on a green cell listing table, the system colors the cell gray to say the slot is unlisted  
**And** when I confirm my new listing time slots and enter a new Price/Hour value and click “Update”, the system asks for confirmation with a small pop up message saying “Are you sure you wish to update your listing?” and "Yes" and "No" buttons.  
&emsp;**And** if I click “Yes”, **Then** the system submits my update parking spot details request
&emsp;**And** the "Price/Hour" fields are updated with the new information,  
&emsp;**And** the large text that says “Unlisted” is replaced by the table with information about the time slot that I have entered  
&emsp;**And** the system displays a flash message saying “Parking spot listed successfully”  
&emsp;**And** when I open the Listing Page and search for my parking spot, the system shows my parking spot listing”  
&emsp;**Or** if I click “No”, the small pop up message will disappear and the system displays the listing table and the “Price/hour” field in the same state before I pressed the “Cancel” button  
**Or** if I click on the "Cancel" button, the system displays a message “Are you sure you wish to cancel without saving any changes?” and "Yes" and "No" buttons  
&emsp;**And** if I click "Yes", the system shows me the current Parking Spot page, with the large text saying “Unlisted”  
&emsp;**Or** if I click “No”, the small pop up message will disappear and the system displays the listing table and the Price/hour field in the same state as before I pressed the “Cancel” button”  

- **User Story 4:** As a seller, I want to remove my parking spot listing from the marketplace if the parking spot is no longer available.
	- **Scenario:** 
“Given that I am a user who owns parking spots, when I open my Personal Parking Spots Page, the system shows a list of my parking spots  
**When** I click on one, the system shows the information for that parking spot and a button that says "Unlist"  
**And** when I click the "Unlist" button, the system asks for confirmation with a small pop-up message saying “Any active listings of the parking spot will be removed.  Are you sure you wish to remove your parking spot?” and "Yes" and "No" buttons  
**And** if I click "Yes", the system displays a "Parking spot unlisted" flash message
**Then** the system takes me back to that Parking Spot page and I can see the parking spot is not listed anymore.  
**Or** if I click “No”, the small pop-up message will disappear and I will remain on that Parking Spot page in the same state as before I pressed the “Unlist” button”  

### Core Feature 2: Booking
**Bookers should be able to book parking spot listings and review their reservations with ease.**

- **User Story 1:** As a booker, I want to book a parking spot I like so I can fulfill my parking needs.
	- **Scenario:** 
“Given that I am a user with booking intentions  
**When** I open the Listings page, the system shows me a list of parking spots with a name and photo for each  
**And** when I click on one, the system shows me the address, price/hour, availability, shelter, plug-in, and charging station information along with a picture of the parking spot and "Book" button  
**And** when I click the “Book” button, the system takes me to a Confirm Booking page with a “Time Period”, “License Plate” fields and “Pay” button  
**And** when I fill the “Time Period” field with a valid time period  
**And** I fill the “License Plate” field with a valid license plate  
**And** I click the “Pay” button  
**Then** the system takes me to the Payment page with “Credit Card Number”, “CVV”, “Expiry Date” fields and a "Confirm" button  
**And** when I enter valid card payment details in the above fields and click the "Confirm" button  
**Then** the system processes my payment and if successful  
**Then** the system displays a flash message saying "Booking successful" and takes me back to the Listings page  
**And** when I open the Booking History page, the system shows this booking in the list”

- **User Story 2:** As a booker, I want to see my active bookings so I can review my parking spot reservations.
	- **Scenario:** 
“Given that I am a user with bookings  
**When** I open the Booking History page, the system shows me a list of my bookings  
**And** when I click on a booking, the system shows me information about the booking, including shelter, plug-in, charge station information, address, amount paid, and time period of the booking”  
### Core Feature 3: Viewing Parking Spots
Bookers should be able to quickly browse through a list of parking spots by harnessing filter, sort, and search aids to find the most convenient spot.

- **User Story 1:** As a booker, I want to see all available parking spots in an organized way so I can easily decide which one to book.
	- **Scenario:** 
“Given that I am a user with booking intentions  
**When** I open the Listings page, the system shows me a list of parking spots with a name and photo for each  
**Where** the default order of the list is by descending price/hour and distance from my geographical location”  

- **User Story 2:** As a booker, I want to search for parking spots near my destination so I can quickly narrow down the options.
	- **Scenario:** 
“Given that I am a user with booking intentions  
**When** I open the Listings page, the system shows me a list of parking spots with a name and photo for each and a search bar with “Address” and “Distance From Address” fields  
**And** when I enter a valid address in the “Address” field   
**And** when I enter a valid distance in the “Distance From Address” field  
**Then** the system updates the Listings page with listings that match the inputted search criteria”  

- **User Story 3:** As a booker, I want to see the available parking spots arranged by criteria such as price and distance from my destination so I can choose the best option.
	- **Scenario:**
“Given that I am a user with booking intentions  
**When** I open the Listings page, the system shows me a list of parking spots with a name and photo for each and a “Sort By” button  
**And** when I click on the “Sort By” button  
**Then** the system will show a small window with four options: a “Price/Hour (ascending)”, “Price/Hour (descending)”, “Distance (ascending)”, and “Distance (descending)” button.  
**And** if I click "Price/Hour (ascending)" button the system sorts the parking spots by price/hour in ascending order
**Or** if I click "Price/Hour (descending)" button the system sorts the parking spots by price/hour in descending order  
**Or** if I click “Distance (ascending)" button the system sorts the parking spots by distance in ascending order relative to the address if an address was already inputted in the search feature or from my geographical location otherwise   
**Or** if I click “Distance (descending)" button the system sorts the parking spots by distance in descending order relative to the address if an address was already inputted in the search feature or from my geographical location otherwise  
**Then** the system updates the Listings page with listings that match the inputted sorting criteria”    

- **User Story 4:** As a booker, I want to filter out available parking spots by criteria such as having shelter or having charging stations because I need to protect my car from the weather or require electricity.
	- **Scenario:** 
“Given that I am a user with booking intentions  
**When** I open the Listings page, the system shows me a list of parking spots with a name and photo for each and a “Filter” button  
**And** when I click on the “Filter” button, the system shows me a small window with three options: a “Plug-in”, "Charging station", and "Shelter" checkboxes along with an “Ok” button  
**And** when I check a combination of the above options and press the “Ok” button
**Then** the system updates the Listings page with listings that match the inputted filter criteria”    

- **User Story 5:** As a booker, I want to save my destination and preferences to get recommendations so I can save time searching for a parking spot.
	- **Scenario:** 
“Given that I am a user with booking intentions  
**When** I go to the Listings page, the system shows me a list of available parking spots and a button that says "Auto-Recommend"   
**And** when I click the “Auto-Recommend” button, the system shows me a Preferences form with “Address”, “Distance From Address”, “Price/hour”, “Time Period” fields and “Plug-in”, "Charging station", and "Shelter" checkboxes along with an “Ok” button  
**Where** the state of the fields and checkboxes match what was entered by the user since the last time the Preferences form was modified and the “Ok” button was pressed
**And** when I fill in all the above fields with valid values and check the boxes for and necessities I need  
**And** I press the “Ok” button  
**Then** the system shows updates the Listings page with listings that meet my saved preferences”  

### Core Feature 4: Parking Spot Management
**Sellers should be able add and remove parking spots from their user profile, track how their parking spots are performing and view a list of booking for each parking spot.**

- **User Story 1:** As a seller, I want to see all the parking spots I am offering in an organized way so I can easily find the listing I am looking for.
	- **Scenario:** 
“Given that I am a user with selling intentions  
**When** I open my Personal Parking Spots page, the system shows a list of my parking spots  
**And** there is a tag on each parking spot that either says “listed” and/or “auctioned” for a parking spot that is actively listed and/or has auctioned time slots
So I know which parking spot is being listed, auctioned, or both at a glance”  

- **User Story 2:** As a seller, I want to see a summary of how my parking spots are performing so I can make more educated decisions.
	- **Scenario:** 
“Given that I am a user with selling intentions  
**When** I open my Analytics page, the system shows me information including the total  revenue and total time used for each of my parking spots ordered by total revenue and total time used in descending order and all-time total revenue earned and total time used across all my parking spots  
**And** when I click on any of the parking spots, it will bring me to the Parking Spot page of that parking spot”

- **User Story 3:** As a seller, I want to see a summary of all the bookings I have for a given parking spot so I can verify if a parked car is authorized.
	- **Scenario:** 
“Given that I am a user with selling intentions  
**When** I open my Personal Parking Spots page, the system shows a list of my parking spots  
**And** when I click on one, I am presented a Parking Spot information screen which displays a listing table  
**Where** the listing table’s column is 7 week days from Monday to Sunday and the column is time from 7:00 am to 12:00 pm, each row is 1 hour; green cell means a listing slot, gray cell means an unlisted slot, and blue means a booked slot  
**And** when I click on a blue cell, the system generates a pop-up window that contains the details of the booker who paid for the booking that includes the selected time slot”  

- **User Story 4:** As a seller, I want to add my vacant parking spots with all the relevant information to my user profile so I can later create a listing on a marketplace for bookers can consider booking it.
	- **Scenario:** 
“Given that I am a user with selling intentions  
**When** I open my Personal Parking Spots page, the system shows a list of my parking spots and a button that says "Add"  
**When** I click the "Add" button, the system takes me to an Add New Parking Spot form with an “Address” field, an element to upload a photo, “Plug-in”, "Charging station", and "Shelter" checkboxes along with a button that says "Add”  
**When** I fill in the above fields and check the appropriate checkboxes, and upload a photo and click the "Add" button  
**Then** the system adds my new parking spot to my user profile  
**And** the system shows an "Parking spot added successfully" flash message and takes me to my Personal Parking Spots page with the newly added parking spot in the list”  

- **User Story 5:** As a seller, I want to quickly remove my parking spots from my user profile that are no longer available.
	- **Scenario:** 
“Given that I am a user with saved parking spots, when I open my Personal Parking Spots page, the system shows a list of my parking spots  
**When** I click on one, the system shows a Parking Spot page with the information for that parking spot and a button that says "Remove".  
**And** when I click the "Remove" button, the system asks for confirmation with a small pop up message saying “Any active listings of the parking spot will also be removed.  Are you sure you wish to remove your parking spot?” and "Yes" and "No" buttons  
**And** if I click "Yes", the system displays a "Parking spot removed" flash message
**Then** the system removes my parking spot and any corresponding listings and takes me back to my Personal Parking Spots page  
**Or** if I click “No”, the small pop up message will disappear and I will remain on the Edit Parking Spot page in the same state before I pressed the “Remove” button”  


### Core Feature 5: Map and Routing
**Bookers should be able to digitally find parking spots through the aid of an interactive map and physically arrive at the parking spot through routing.**

- **User Story 1:** As a booker, I want to see the locations of all available parking spots on a map so I can easily compare their distances from my destination.
	- **Scenario:** 
“Given that I am a user with booking intentions  
**When** I open the Listings page, the system shows me a list of parking spots with a name and photo for each and a "Map View" button  
**And** when I click the "Map View" button  
**Then** the system brings me to a Map View screen, showing me a map with the parking spots marked on it”  

- **User Story 2:** As a booker, I want to get directions to the parking spot I booked so I can easily navigate there.
	- **Scenario:** 
“Given that I am a user with booking intentions  
**When** I open my Booking History page, the system shows me a list of my bookings  
**And** when I click on a booking, the system shows me the information of the booking including shelter, plug-in, charge station information, address, amount paid, and time period of the booking along with a "Directions" button  
**And** when I click the "Directions" button, the system opens a map application with directions to the parking spot from the user’s current location”  

### Core Feature 6: Bidding
**Sellers should be able to auction their high demand parking spots to maximize profits and buyers should be able to bid on such auctions to compete for ideal parking.**

- **User Story 1:** As a seller, I want to offer my high demand parking spots in the form of an auction instead of fixed pricing to maximize my profits.
	- **Scenario:** 
“Given that I am a user with selling intentions, when I open my Personal Parking Spots page, the system shows me a list of my parking spots  
**And** when I click on a parking spot, the system shows a Parking Spot Details page with the information for that parking spot including an "Auction" button  
**And** when I click on the "Auction" button, the system shows me a form containing "Time Period" and "Starting Price" fields and “Submit” and “Cancel” buttons  
**And** when I fill in valid information in the above fields and click the "Submit" button  
**Then** the system asks for confirmation with a small pop up message saying “Are you sure you wish to auction your parking spot?” and "Yes" and "No" buttons  
&emsp;**And** if I click “Yes”, **Then** the system submits my auction request  
&emsp;**And** the system displays a flash message saying “Parking spot auctioned successfully”  
&emsp;**And** when I open the Listings page and search for my parking spot, the system shows my parking spot auction  
&emsp;**Or** if I click “No”, the small pop up message will disappear and the system will return my to the Parking Spot Details page  
**Or** if I click on the "Cancel" button, the system displays a message “Are you sure you wish to cancel without saving any changes?” and "Yes" and "No" buttons  
&emsp;**And** if I click "Yes", the system shows me the current Parking Spot Details page in the same state before I pressed the “Auction” button  
&emsp;**Or** if I click “No”, the small pop up message will disappear and the system displays the values in the “Time Period” and “Starting Price” fields in the same state before I pressed the “Cancel” button”  

- **User Story 2:** As a booker, I want to place a bid on parking spots offered as an auction because I am willing to compete with others for a convenient location.
	- **Scenario:** 
“Given that I am a user with booking intentions  
**When** I open the Listings page, the system shows me a list of available parking spots along with a photo and a name for each  
**And** when I click on one, the system shows me the address, price/hour, availability, shelter, plug-in, and charging station information along with a picture of the parking spot and an “Auctions” tab  
**And** when I click on the “Auctions” tab, the system takes me to an Available Auctions page displaying the available time periods, the current bid price, and a “Bid” button  
**And** when I click the “Bid” button, then the system displays a pop-up window saying “Are you sure you wish to increase the current bid price by $0.25?” and “Yes” and “No” buttons  
**And** if I click “Yes”, the system automatically increases the current bid price by $0.25 and closes the pop-up window and the system displays a "Bid placed successfully" flash message
**Or** if I click “No”, then the pop-up window closes”

- **User Story 3:** As a booker, I want to be able to automatically bid up on certain parking spot auctions so I don’t have to keep coming back and checking the auction.
	- **Scenario:** Given that I am a user with booking intentions
**When** I open the Listings page, the system shows me a list of available parking spots alone with a photo and name for each  
**And** when I click on one, the system shows me the address, price/hour, availability, shelter, plug-in, and charging station information along with a picture of the parking spot and an “Auctions” tab  
**And** when I click on the “Auctions” tab, the system takes me to an Available Auctions page displaying the available time periods, the current bid price, and a “Bid” button along with an “Auto-bid” checkbox  
**And** if I click on the “Auto-bid” checkbox, the system will display a pop-up window with an “Enter maximum bid price” field and “Ok” and “Cancel” buttons  
**And** when I enter a valid value in the “Enter maximum bid price” field and click the “Yes” button then the system will enable automatic bidding to up-bid by $0.25 whenever a different user out-bids my current bid price up to the entered maximum bid price  
**Or** if I click the “Cancel” button then the pop-up window will close”  

# Non-functional Manual Tests

### Android

### Log In, Sign Up, Log Out

**Test 1**
Steps:
1. From the Login Screen, press the "Register Now" link
2. Select the "Name" text field, input "John Appleseed"
3. Select the "Email" text field, input "user@example.com"
4. Select the "Password" text field, input "password"
5. Select the "Confirm Password" text field, input "password"
6. Press the "Register" button
*Expected: "Registered successfully" is displayed and app navigates away from the Login Screen*

**Test 2**
1. From the Profile Screen, press the "Logout" button
*Expected: App navigates to the Login Screen*

**Test 3**
1. From the Login Screen, select the "email" text field, input "user@example.com"
2. Select the "password" text field, input "password"
3. Press the "Login" button
*Expected: "Logged in successfully" is displayed and app navigates away from the Login Screen*

### Webapp

### Log In, Sign Up, Log Out

**Test 1**
1. Open your browser of choice, navigate to "http://localhost:5173/auth/login"
2. Press the "Create one" link
3. Select the "First Name" text field, input "Robert"
4. Select the "Last Name" text field, input "Guderian"
5. Select the "Email" text field, input "robg@cs.umanitoba.ca"
6. Select the "Password" text field, input "password"
7. Select the "Confirm Password" text field, input "password"
8. Press the "Sign Up" button
*Expected: "Account created successfully" is displayed*

**Test 2**
1. Open your browser of choice, navigate to "http://localhost:5173/auth/login"
2. Select the "Email" text field, input "robg@cs.umanitoba.ca"
3. Select the "Password" text field, input "password"
4. Press the "Login" button
*Expected: Website redirects to a different page*
5. Press the "my profile" button
*Expected: "Robert Guderian" and "robg@cs.umanitoba.ca" are displayed*
6. Press the "Logout" button
*Expected: Website redirects to a the login page"
