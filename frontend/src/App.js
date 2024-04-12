import { HiShoppingCart, HiXCircle } from "react-icons/hi";
import React, { useState, useEffect } from 'react';
import { Button, Card, Table } from "flowbite-react";
import './App.css';
import Products from './products.json';

function Product({ product, addToCart }) {
  const handleAddToCart = () => {
    addToCart(product);
  };

  return (
    <Card className="max-w-sm mr-1.5 mb-1.5 float-left">
      <h5 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
        {product.name}
      </h5>
      <div className="flex items-center justify-between">
        <span className="text-3xl font-bold text-gray-900 dark:text-white">{product.price} €</span>
        <Button onClick={handleAddToCart}><HiShoppingCart className="h-5 w-5" /></Button>
      </div>
    </Card>
  );
}

function ProductList({ products, addToCart }) {
  return (
    <div className="grow">
      {products.map(product => (
        <Product key={product.id} product={product} addToCart={addToCart} />
      ))}
    </div>
  );
}

function Cart({ cart, cartValue }) {
  return (
    <div className="w-30">
      <Table striped>
        <Table.Head>
          <Table.HeadCell>Product</Table.HeadCell>
          <Table.HeadCell>Quantity</Table.HeadCell>
          <Table.HeadCell>Total Price</Table.HeadCell>
          <Table.HeadCell>Remove</Table.HeadCell>
        </Table.Head>
        <Table.Body>
          {cart.map(cartElement => (
            <Table.Row key={cartElement.id}>
              <Table.Cell>{cartElement.name}</Table.Cell>
              <Table.Cell>{cartElement.count}</Table.Cell>
              <Table.Cell>{cartElement.totalPrice} €</Table.Cell>
              <Table.Cell><Button><HiXCircle /></Button></Table.Cell>
            </Table.Row>
          ))}
          <Table.Row>
            <Table.Cell>Total</Table.Cell>
            <Table.Cell></Table.Cell>
            <Table.Cell>{cartValue} €</Table.Cell>
            <Table.Cell></Table.Cell>
          </Table.Row>
        </Table.Body>
      </Table>
    </div>
  );
}

function App() {
  const [cart, setCart] = useState([]);
  
  const addToCart = (product) => {
    console.log("wuff" + product);
    const existingProductIndex = cart.findIndex(item => item.id === product.id);
    if (existingProductIndex !== -1) {
      const updatedCart = [...cart];
      updatedCart[existingProductIndex].count++;
      updatedCart[existingProductIndex].totalPrice = updatedCart[existingProductIndex].count * updatedCart[existingProductIndex].price;
      setCart(updatedCart);
    } else {
      const updatedProduct = { ...product, count: 1, totalPrice: product.price };
      setCart([...cart, updatedProduct]);
    }
  };

  return (
    <div className="App">
      <div className="flex">
        <ProductList products={Products} addToCart={addToCart} />
        <Cart cart={cart} cartValue={0} />
      </div>
    </div>
  );
}

export default App;
