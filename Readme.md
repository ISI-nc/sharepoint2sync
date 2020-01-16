# sharepoint2sync

Synchronize sharepoint atom feed endpoint with a sync2kafka server.


## Spec

This app woks with shareepoint api endpoints. Specifically with those returning arrays of values.

Sample request response body from sharepoint :
```json
{
  "value": [
    {
      "FileSystemObjectType": 0,
      "Id": 4,
      "ID": 4,
      "ContentTypeId": "0x0100BE43D108C2F78345AE6B759AF12B2D11",
      "Title": "janv-13",
      "Modified": "2019-04-18T20:48:47Z",
      "Created": "2018-11-21T21:52:44Z",
      "AuthorId": 1,
      "EditorId": 13,
      "OData__UIVersionString": "1.0",
      "Attachments": false,
      "GUID": "e8067b7b-6daa-4b8c-98e1-930ac4c0e461",
      "Colonne1": "10",
      "Colonne2": "11",
      "Colonne3": "12",
      "Colonne4": "13",
      "Colonne5": "14",
      "Colonne8": "15",
      "Colonne9": "16",
      "Colonne10": null,
      "ferme_x0020__x00e9_olienne_x0020": "540",
      "GE_x0020_EEC": "985"
    },
    {
      "FileSystemObjectType": 0,
      "Id": 5,
      "ID": 5,
      "ContentTypeId": "0x0100BE43D108C2F78345AE6B759AF12B2D11",
      "Title": "f\u00e9vr-13",
      "Modified": "2019-04-18T20:53:23Z",
      "Created": "2018-11-21T21:52:44Z",
      "AuthorId": 1,
      "EditorId": 13,
      "OData__UIVersionString": "1.0",
      "Attachments": false,
      "GUID": "ea45724a-4595-40a5-bf74-0068bf909e18",
      "Colonne1": "1",
      "Colonne2": "2",
      "Colonne3": "3",
      "Colonne4": "4",
      "Colonne5": "5",
      "Colonne8": "6",
      "Colonne9": "7",
      "Colonne10": null,
      "ferme_x0020__x00e9_olienne_x0020": "786",
      "GE_x0020_EEC": "884"
    } 
  ]
}
```

generated messages :

```json
4 -> {
            "FileSystemObjectType": 0,
            "Id": 4,
            "ID": 4,
            "ContentTypeId": "0x0100BE43D108C2F78345AE6B759AF12B2D11",
            "Title": "janv-13",
            "Modified": "2019-04-18T20:48:47Z",
            "Created": "2018-11-21T21:52:44Z",
            "AuthorId": 1,
            "EditorId": 13,
            "OData__UIVersionString": "1.0",
            "Attachments": false,
            "GUID": "e8067b7b-6daa-4b8c-98e1-930ac4c0e461",
            "Colonne1": "10",
            "Colonne2": "11",
            "Colonne3": "12",
            "Colonne4": "13",
            "Colonne5": "14",
            "Colonne8": "15",
            "Colonne9": "16",
            "Colonne10": null,
            "ferme_x0020__x00e9_olienne_x0020": "540",
            "GE_x0020_EEC": "985"
          }
```

```json
5 -> {
           "FileSystemObjectType": 0,
           "Id": 5,
           "ID": 5,
           "ContentTypeId": "0x0100BE43D108C2F78345AE6B759AF12B2D11",
           "Title": "f\u00e9vr-13",
           "Modified": "2019-04-18T20:53:23Z",
           "Created": "2018-11-21T21:52:44Z",
           "AuthorId": 1,
           "EditorId": 13,
           "OData__UIVersionString": "1.0",
           "Attachments": false,
           "GUID": "ea45724a-4595-40a5-bf74-0068bf909e18",
           "Colonne1": "1",
           "Colonne2": "2",
           "Colonne3": "3",
           "Colonne4": "4",
           "Colonne5": "5",
           "Colonne8": "6",
           "Colonne9": "7",
           "Colonne10": null,
           "ferme_x0020__x00e9_olienne_x0020": "786",
           "GE_x0020_EEC": "884"
         }
```

