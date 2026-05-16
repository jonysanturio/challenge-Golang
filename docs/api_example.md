# API Examples

## Products

### Get All Products
**Request**
```http
GET /api/products?page=1&limit=10&sort_by=price&sort_order=asc
 ```

 **Response**
 ```json
 {
   "data": [
     {
       "id": 1,
       "name": "Smartphone X",
       "description": "Latest flagship smartphone",
       "price": 999.99,
       "stock": 50,
       "created_at": "2026-05-15T10:30:00Z",
       "updated_at": "2026-05-15T10:30:00Z",
       "categories": [
         {
           "id": 1,
           "name": "Electronics",
           "description": "Electronic devices and gadgets"
         }
       ]
     }
   ],
   "pagination": {
     "page": 1,
     "limit": 10,
     "total": 1,
     "pages": 1
   }
 }
 ```

 ### Get Product by ID
 **Request**
 ```http
 GET /api/products/1
 ```

 **Response**
 ```json
 {
   "data": {
     "id": 1,
     "name": "Smartphone X",
     "description": "Latest flagship smartphone",
     "price": 999.99,
     "stock": 50,
     "created_at": "2026-05-15T10:30:00Z",
     "updated_at": "2026-05-15T10:30:00Z",
     "categories": [
       {
         "id": 1,
         "name": "Electronics",
         "description": "Electronic devices and gadgets"
       }
     ]
   }
 }
 ```

 ### Create Product
 **Request**
 ```http
 POST /api/products
 Authorization: Bearer <jwt-token>
 Content-Type: application/json

 {
   "name": "Wireless Earbuds",
   "description": "True wireless stereo earbuds",
   "price": 149.99,
   "stock": 75,
   "category_ids": [1]
 }
 ```

 **Response (201 Created)**
 ```json
 {
   "data": {
     "id": 2,
     "name": "Wireless Earbuds",
     "description": "True wireless stereo earbuds",
     "price": 149.99,
     "stock": 75,
     "created_at": "2026-05-15T14:22:00Z",
     "updated_at": "2026-05-15T14:22:00Z",
     "categories": [
       {
         "id": 1,
         "name": "Electronics",
         "description": "Electronic devices and gadgets"
       }
     ]
   }
 }
 ```

 ### Update Product
 **Request**
 ```http
 PUT /api/products/2
 Authorization: Bearer <jwt-token>
 Content-Type: application/json

 {
   "price": 129.99,
   "stock": 80
 }
 ```

 **Response**
 ```json
 {
   "data": {
     "id": 2,
     "name": "Wireless Earbuds",
     "description": "True wireless stereo earbuds",
     "price": 129.99,
     "stock": 80,
     "created_at": "2026-05-15T14:22:00Z",
     "updated_at": "2026-05-15T15:45:00Z",
     "categories": [
       {
         "id": 1,
         "name": "Electronics",
         "description": "Electronic devices and gadgets"
       }
     ]
   }
 }
 ```

 ### Delete Product
 **Request**
 ```http
 DELETE /api/products/2
 Authorization: Bearer <jwt-token>
 ```

 **Response (200 OK)**
 ```json
 {
   "message": "Product deleted successfully"
 }
 ```

 ### Get Product History
 **Request**
 ```http
 GET /api/products/1/history?start=2026-05-01&end=2026-05-31
 ```

 **Response**
 ```json
 {
   "data": [
     {
       "id": 1,
       "product_id": 1,
       "price": 999.99,
       "stock": 50,
       "changed_at": "2026-05-15T10:30:00Z"
     },
     {
       "id": 2,
       "product_id": 1,
       "price": 899.99,
       "stock": 45,
       "changed_at": "2026-05-16T09:15:00Z"
     }
   ]
 }
 ```

 ## Categories

 ### Get All Categories
 **Request**
 ```http
 GET /api/categories?page=1&limit=5
 ```

 **Response**
 ```json
 {
   "data": [
     {
       "id": 1,
       "name": "Electronics",
       "description": "Electronic devices and gadgets",
       "created_at": "2026-05-15T10:30:00Z",
       "updated_at": "2026-05-15T10:30:00Z",
       "products": [
         {
           "id": 1,
           "name": "Smartphone X",
           "price": 999.99,
           "stock": 50
         }
       ]
     }
   ],
   "pagination": {
     "page": 1,
     "limit": 5,
     "total": 1,
     "pages": 1
   }
 }
 ```

 ### Create Category
 **Request**
 ```http
 POST /api/categories
 Authorization: Bearer <jwt-token>
 Content-Type: application/json

 {
   "name": "Accessories",
   "description": "Electronic accessories and peripherals"
 }
 ```

 **Response (201 Created)**
 ```json
 {
   "data": {
     "id": 2,
     "name": "Accessories",
     "description": "Electronic accessories and peripherals",
     "created_at": "2026-05-15T16:30:00Z",
     "updated_at": "2026-05-15T16:30:00Z",
     "products": []
   }
 }
 ```

 ## Search

 ### Search Products
 **Request**
 ```http
 GET /api/search?q=phone&type=product&page=1&limit=10
 ```

 **Response**
 ```json
 {
   "data": [
     {
       "id": 1,
       "name": "Smartphone X",
       "description": "Latest flagship smartphone",
       "price": 999.99,
       "stock": 50,
       "created_at": "2026-05-15T10:30:00Z",
       "updated_at": "2026-05-15T10:30:00Z",
       "categories": [
         {
           "id": 1,
           "name": "Electronics",
           "description": "Electronic devices and gadgets"
         }
       ]
     }
   ],
   "pagination": {
     "page": 1,
     "limit": 10,
     "total": 1,
     "pages": 1
   },
   "search": {
     "type": "product",
     "query": "phone"
   }
 }
 ```

 ### Search Categories
 **Request**
 ```http
 GET /api/search?q=electronics&type=category
 ```

 **Response**
 ```json
 {
   "data": [
     {
       "id": 1,
       "name": "Electronics",
       "description": "Electronic devices and gadgets",
       "created_at": "2026-05-15T10:30:00Z",
       "updated_at": "2026-05-15T10:30:00Z",
       "products": [
         {
           "id": 1,
           "name": "Smartphone X",
           "price": 999.99,
           "stock": 50
         }
       ]
     }
   ],
   "pagination": {
     "page": 1,
     "limit": 10,
     "total": 1,
     "pages": 1
   },
   "search": {
     "type": "category",
     "query": "electronics"
   }
 }
 ```

 ## WebSocket Examples

 ### Connecting to WebSocket
 ```javascript
 // JavaScript example
 const socket = new WebSocket('ws://localhost:8080/ws');
 socket.onopen = function(event) {
   console.log('Connected to WebSocket server');
   // Send a message if needed
 // socket.send(JSON.stringify({ action: 'subscribe', channel: 'products' }));
};

socket.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('Received update:', data);

  // Handle different event types
  switch(data.event) {
    case 'product.created':
      // Handle new product
      break;
    case 'product.updated':
      // Handle product update
      break;
    case 'product.deleted':
      // Handle product deletion
      break;
    case 'category.created':
      // Handle new category
      break;
    case 'category.updated':
      // Handle category update
      break;
    case 'category.deleted':
      // Handle category deletion
      break;
  }
};

socket.onclose = function(event) {
  console.log('Disconnected from WebSocket server');
};

socket.onerror = function(error) {
  console.error('WebSocket error:', error);
};
```

#### Expected WebSocket Message Format
When a product is created:
```json 
{
  "event": "product.created",
  "data": {
    "id": 3,
    "name": "New Product",
    "description": "Product description",
    "price": 29.99,
    "stock": 100,
    "created_at": "2026-05-16T10:00:00Z",
    "updated_at": "2026-05-16T10:00:00Z",
    "categories": [
      {
        "id": 2,
        "name": "Accessories",
        "description": "Electronic accessories and peripherals"
       }
     ]
   }
 }
 ```

 When a product is updated:
 ```json
 {
   "event": "product.updated",
   "data": {
     "id": 1,
     "name": "Smartphone X",
     "description": "Latest flagship smartphone",
     "price": 899.99,
     "stock": 45,
     "created_at": "2026-05-15T10:30:00Z",
     "updated_at": "2026-05-16T11:30:00Z",
     "categories": [
       {
         "id": 1,
         "name": "Electronics",
         "description": "Electronic devices and gadgets"
       }
     ]
   }
 }
 ```

 When a product is deleted:
 ```json
 {
   "event": "product.deleted",
   "data": {
     "id": 1
   }
 }
 ```

 Similar format applies for category events.