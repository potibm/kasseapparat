import { HiXCircle } from 'react-icons/hi'
import { Button, Table } from 'flowbite-react'

export default function Cart ({ cart, removeFromCart, removeAllFromCart, checkoutCart, currency }) {
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
              <Table.Row key={cartElement.ID}>
                <Table.Cell className="whitespace-nowrap">{cartElement.Name}</Table.Cell>
                <Table.Cell className="text-right">{cartElement.count}</Table.Cell>
                <Table.Cell className="text-right">{currency.format(cartElement.totalPrice)}</Table.Cell>
                <Table.Cell><Button color="failure" onClick={() => removeFromCart(cartElement)}><HiXCircle /></Button></Table.Cell>
              </Table.Row>
            ))}
            <Table.Row>
              <Table.Cell className="uppercase font-bold">Total</Table.Cell>
              <Table.Cell></Table.Cell>
              <Table.Cell className="font-bold text-right">{currency.format(cart.reduce((total, item) => total + item.totalPrice, 0))}</Table.Cell>
              <Table.Cell>{cart.length
                ? (
                <Button color="failure" onClick={() => removeAllFromCart()}><HiXCircle /></Button>
                  )
                : (
                <Button disabled color="failure"><HiXCircle /></Button>
                  )}</Table.Cell>
            </Table.Row>
          </Table.Body>
        </Table>

        <Button {...(cart.length === 0 && { disabled: true })} color="success" className="w-full mt-2 uppercase" onClick={checkoutCart}>
          Checkout&nbsp;
          {cart.length > 0 && currency.format(cart.reduce((total, item) => total + item.totalPrice, 0))}
        </Button>
      </div>
  )
}
