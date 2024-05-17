import { HiXCircle } from 'react-icons/hi'
import React, { useState, useEffect } from 'react'
import { Button, Table } from 'flowbite-react'
import { fetchPurchases } from '../hooks/Api'

export default function PurchaseHistory ({history, currency}) {

    return (
      <div className='mt-10'>
        <Table striped>
          <Table.Head>
            <Table.HeadCell>Date</Table.HeadCell>
            <Table.HeadCell className="text-right">Total Price</Table.HeadCell>
            <Table.HeadCell>Remove</Table.HeadCell>
          </Table.Head>
          <Table.Body>
          {history.slice(0, 2).map(purchase => (
              <Table.Row key={purchase.id}>
                <Table.Cell className="text-right">{new Date(purchase.createdAt).toLocaleString('de-DE', { weekday: 'long', hour: '2-digit', minute: '2-digit' })}</Table.Cell>
                <Table.Cell className="text-right">{currency.format(purchase.totalPrice)}</Table.Cell>
                <Table.Cell><Button><HiXCircle /></Button></Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </div>
    )
}
