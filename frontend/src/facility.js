export default {
  "timezone": {
    "name": "America/Anchorage",
    "offset": -9,
    "dst": -8,
  },
  "views": [
    { "name": "Allen AAF ATCT", "facilities": ["BIG"] },
    { "name": "Anchorage ATCT", "facilities": ["ANC"] },
    { "name": "Bethel ATCT", "facilities": ["BET"] },
    { "name": "Eielson AFB ATCT", "facilities": ["EIL"] },
    { "name": "Fairbanks ATCT", "facilities": ["FAI"] },
    { "name": "Juneau ATCT", "facilities": ["JNU"] },
    { "name": "Kenai ATCT", "facilities": ["ENA"] },
    { "name": "King Salmon ATCT", "facilities": ["AKN"] },
    { "name": "Kodiak ATCT", "facilities": ["ADQ"] },
    { "name": "Ladd AAF ATCT", "facilities": ["FBK"] },
    { "name": "Lake Hood ATCT", "facilities": ["LHD"] },
    { "name": "Merrill Field ATCT", "facilities": ["MRI"] },
    { "name": "Fairbanks TRACON", "facilities": ["FAI", "FBK", "EIL"] },
    { "name": "Anchorage TRACON", "facilities": ["ANC", "LHD", "MRI", "EDF"] },
    { "name": "ZAN North Area", "facilities": ["FAI", "AKN", "BET", "FBK", "EIL", "BIG"] },
    { "name": "ZAN South Area", "facilities": ["ANC", "JNU", "ADQ", "ENA", "MRI", "EDF"] },
  ],
  "airports": [
    { "name": "ANC", "hours": { "continuous": true } },
    {
      "name": "ADQ", "hours": {
        "continuous": false, "schedule": [
          { "start": { "month": 4, "day": 1 }, "end": { "month": 9, "day": 30 }, "open": "07:00", "close": "22:00", "local": true },
          { "start": { "month": 10, "day": 1 }, "end": { "month": 3, "day": 31 }, "open": "06:30", "close": "20:00", "local": true },
        ]
      }
    },
    {
      "name": "BET", "hours": {
        "continuous": false, "schedule": [
          { "start": { "month": 4, "day": 1 }, "end": { "month": 10, "day": 31 }, "open": "07:00", "close": "22:00", "local": true },
          { "start": { "month": 11, "day": 1 }, "end": { "month": 3, "day": 31 }, "open": "07:00", "close": "20:00", "local": true },
        ]
      }
    },
    {
      "name": "BIG", "hours": {
        "continuous": false, "schedule": [
          { "whenDST": true, "open": "16:15", "close": "00:00", "days": [1, 2, 3, 4, 5], "local": false },
          { "whenDST": false, "open": "17:15", "close": "01:00", "days": [1, 2, 3, 4, 5], "local": false },
        ]
      }
    },
    { "name": "EDF", "hours": { "continuous": true } },
    {
      "name": "EIL", "hours": {
        "continuous": false, "schedule": [
          { "whenDST": true, "open": "15:00", "close": "07:00", "local": false },
          { "whenDST": false, "open": "16:00", "close": "08:00", "local": false },
        ]
      }
    },
    {
      "name": "ENA", "hours": {
        "continuous": false, "schedule": [
          { "start": { "month": 5, "day": 1 }, "end": { "month": 9, "day": 30 }, "open": "06:00", "close": "22:00", "local": true },
          { "start": { "month": 10, "day": 1 }, "end": { "month": 4, "day": 30 }, "open": "07:00", "close": "21:00", "local": true },
        ]
      }
    },
    { "name": "FAI", "hours": { "continuous": true } },
    {
      "name": "FBK", "hours": {
        "continuous": false, "schedule": [
          { "whenDST": true, "open": "16:00", "close": "07:00", "local": false },
          { "whenDST": false, "open": "17:00", "close": "08:00", "local": false }
        ]
      }
    },
    {
      "name": "JNU", "hours": {
        "continuous": false, "schedule": [
          { "start": { "month": 4, "day": 1 }, "end": { "month": 9, "day": 30 }, "open": "06:00", "close": "23:00", "local": true },
          { "start": { "month": 10, "day": 1 }, "end": { "month": 3, "day": 31 }, "open": "07:00", "close": "20:00", "local": true },
        ]
      }
    },
    { "name": "LHD", "hours": { "continuous": true } },
    {
      "name": "MRI", "hours": {
        "continuous": false, "schedule": [
          { "start": { "month": 1, "day": 1 }, "end": { "month": 12, "day": 31 }, "open": "07:00", "close": "22:00", "local": true },
        ]
      }
    },
  ]
};
