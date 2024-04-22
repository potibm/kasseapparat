import { Button, Card } from 'flowbite-react'
import { HiShoppingCart } from 'react-icons/hi'

export default function Product ({ product, addToCart, currency }) {
  const handleAddToCart = () => {
    addToCart(product)
  }

  return (
      <Card className={`w-3/12 mr-1.5 mb-1.5 float-left ${product.WrapAfter ? 'float-none' : ''}`}>
        <h5 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          {product.Name}
        </h5>
        <div className="flex items-center justify-between">
          <span className="text-3xl font-bold text-gray-900 dark:text-white">{currency.format(product.Price)}</span>
          <Button onClick={handleAddToCart}><HiShoppingCart className="h-5 w-5" /></Button>
        </div>
      </Card>
  )
}
