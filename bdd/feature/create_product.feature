Feature: Create product
  In order to create a product for my restaurant
  As a admin of the fastfoot restaurant
  I need to be able add a product in the restaurant system

  Scenario: then user try to pay the order, success should be displayed
    When I send "POST" request to "/product" with payload:
      """
      {
        "name": "Batata 2",
        "description": "Batata média",
        "category": "Acompanhamento",
        "price": 7.99,
        "images": [
            {
                "imageUrl": "url1"
            }
        ]
      }  
      """
    Then the response code should be 200
    And the response payload should match json:
      """
        {
            "id": 3
        } 
      """