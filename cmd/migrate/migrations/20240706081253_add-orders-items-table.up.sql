CREATE TABLE IF NOT EXISTS order_items (
  id SERIAL NOT NULL,
  orderId INT NOT NULL,
  productId INT NOT NULL,
  quantity INT NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  
  PRIMARY KEY (id),
  FOREIGN KEY (orderId) REFERENCES orders(id),
  FOREIGN KEY (productId) REFERENCES products(id)
);
