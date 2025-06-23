# SumUp Integration

You need a SumUp account and a **SumUp Solo** card reader. Other devices – especially the **Solo Light** – are not supported.

## Configure .env file

To configure SumUp you will need the following information in your `.env` file:

```
PAYMENT_METHODS="CASH,SUMUP"
SUMUP_API_KEY="sup_sk_01234567890abcdef0123456789abcdef"
SUMUP_MERCHANT_CODE="M0123456"
SUMUP_APPLICATION_ID="com.example.kasseapparat"
SUMUP_AFFILIATE_KEY="sup_afk_01234567890abcdef0123456789abcdef"
SUMUP_PUBLIC_URL="https://kasseapparat.example.com"
```

### Merchant Code (Merchant ID)

You can find your Merchant Code on [SumUp Settings](https://me.sumup.com/de-de/settings), displayed directly below your company name.

### API Key

Log in to [SumUp Developer API Keys](https://developer.sumup.com/api-keys) and generate a new key.

Give it a descriptive name and store the key starting with `sup_sk` as `SUMUP_API_KEY` in your .env file.

### Application ID

Log in to [Affiliate Keys](https://developer.sumup.com/affiliate-keys) and add a new Application ID (in the format `com.example.app`).

You can find more details in the [Getting Started guide](https://developer.sumup.com/terminal-payments/introduction/getting-started) of the SumUp SDK/API.

Store this ID as `SUMUP_APPLICATION_ID`.

### Affiliate Key

When creating the Application ID, you should copy the corresponding Affiliate Key from [Affiliate Keys](https://developer.sumup.com/affiliate-keys).

Store it as `SUMUP_AFFILIATE_KEY`. The key starts with `sup_afk`.

### Public URL

If your Kasseapparat setup is publicly accessible via the internet, you can enter its public URL here. Otherwise, leave this field empty.

This URL is used for a **webhook**, which helps improve response times from the terminals.

## Pairing a Reader

Follow the instructions at [Pairing a Solo Reader](https://developer.sumup.com/terminal-payments/guides/pairing-solo) until you see the **pairing code** displayed on the device.

Then go to the **Kasseapparat Admin** and navigate to **SumUp > Readers**.  
Click **"+ Pair"**, enter the pairing code, and provide a descriptive name for the reader.  
A meaningful name is especially useful if you are using multiple readers and need to distinguish between them.

After clicking **"Pair"**, the reader should now be successfully paired.

## Selecting a Reader

Once one or more readers are paired, you need to assign one to each terminal.

In the **SumUp > Readers** section of the admin interface, click **"Use this reader"** next to the reader you want to use.  
The selected reader will be marked as **"Selected"**.

You can now start accepting purchases in the frontend using the selected reader.
