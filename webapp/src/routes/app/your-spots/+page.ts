import type { PageLoad } from "./$types";
const test_data = [{
    id: "1",
    features: {
    "charging_station": true,
    "plug_in": true,
    "shelter": true
  },
  location: {
    "city": "Winnipeg",
    "country_code": "CA",
    "latitude": 1,
    "longitude": 1,
    "postal_code": "RRRRRR",
    "street_address": "230 Portage Avenue"
  },
  isListed: true,
}, {
    id: "2",
    features: {
    "charging_station": true,
    "plug_in": true,
    "shelter": true
  },
  location: {
    "city": "Winnipeg",
    "country_code": "CA",
    "latitude": 1,
    "longitude": 1,
    "postal_code": "RRRRRR",
    "street_address": "230 Portage Avenue"
  },
  isListed: false,
},{
    id: "3",
    features: {
    "charging_station": true,
    "plug_in": true,
    "shelter": true
  },
  location: {
    "city": "Winnipeg",
    "country_code": "CA",
    "latitude": 1,
    "longitude": 1,
    "postal_code": "RRRRRR",
    "street_address": "230 Portage Avenue"
  },
  isListed: true,
}]

export const load: PageLoad = async ({ fetch }) => {
    return {
        spots: test_data,
        hasNext: undefined,
        paging: undefined,
    }
}