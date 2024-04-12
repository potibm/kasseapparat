import { HiShoppingCart, HiXCircle } from "react-icons/hi";
import React, { useState, useEffect } from 'react';
import { Button, Card, Table } from "flowbite-react";
import './App.css';
import Products from './products.json';

let Currency = new Intl.NumberFormat('de-DE', {
  style: 'currency',
  currency: 'EUR',
  minimumFractionDigits: 0,
  maximumFractionDigits: 0,
});

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
        <span className="text-3xl font-bold text-gray-900 dark:text-white">{Currency.format(product.price)}</span>
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

function Cart({ cart, removeFromCart, removeAllFromCart, checkoutCart }) {
  return (
    <div className="w-30">
      <Table striped>
        <Table.Head>
          <Table.HeadCell>Product</Table.HeadCell>
          <Table.HeadCell className="text-right">Quantity</Table.HeadCell>
          <Table.HeadCell className="text-right">Total Price</Table.HeadCell>
          <Table.HeadCell>Remove</Table.HeadCell>
        </Table.Head>
        <Table.Body>
          {cart.map(cartElement => (
            <Table.Row key={cartElement.id}>
              <Table.Cell className="whitespace-nowrap">{cartElement.name}</Table.Cell>
              <Table.Cell className="text-right">{cartElement.count}</Table.Cell>
              <Table.Cell className="text-right">{Currency.format(cartElement.totalPrice)}</Table.Cell>
              <Table.Cell><Button color="failure" onClick={() => removeFromCart(cartElement)}><HiXCircle /></Button></Table.Cell>
            </Table.Row>
          ))}
          <Table.Row>
            <Table.Cell className="uppercase font-bold">Total</Table.Cell>
            <Table.Cell></Table.Cell>
            <Table.Cell className="font-bold text-right">{Currency.format(cart.reduce((total, item) => total + item.totalPrice, 0))}</Table.Cell>
            <Table.Cell>{cart.length ? (
              <Button color="failure" onClick={() => removeAllFromCart()}><HiXCircle /></Button> 
            ) : (
              <Button disabled color="failure"><HiXCircle /></Button> 
            )}</Table.Cell>
          </Table.Row>
        </Table.Body>
      </Table>

      <Button {...(cart.length === 0 && {disabled: true})}  color="success" className="w-full mt-2 uppercase" onClick={checkoutCart}>
        Checkout&nbsp;
        {cart.length && Currency.format(cart.reduce((total, item) => total + item.totalPrice, 0))}
      </Button>
    </div>
  );
}

function App() {
  const [cart, setCart] = useState([]);
  
  const addToCart = (product) => {
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

  const removeFromCart = (product) => {
    const existingProductIndex = cart.findIndex(item => item.id === product.id);
  
    if (existingProductIndex !== -1) {
      setCart([...cart.slice(0, existingProductIndex), ...cart.slice(existingProductIndex + 1)]);
    }
  }

  const removeAllFromCart = () => {
    setCart([]);
  }

  const checkoutCart = () => {
    setCart([]);
  }

  return (
    <div className="App p-2">
      <div className="flex">
        <ProductList products={Products} addToCart={addToCart} />
        <Cart cart={cart} 
          removeFromCart={removeFromCart} 
          removeAllFromCart={removeAllFromCart} 
          checkoutCart={checkoutCart} />
      </div>
    </div>
  );
}

export default App;
