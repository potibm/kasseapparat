# SumUp Integration

You need a SumUp account and a **SumUp Solo** card reader. Other devices – especially the **Solo Light** – are not supported.

## Configuration

Put the SumUp credentials in `config/config.local.yaml`, which is git-ignored:

```yaml
sumup:
  api_key: sup_sk_01234567890abcdef0123456789abcdef
  merchant_code: M0123456
  application_id: com.example.kasseapparat
  affiliate_key: sup_afk_01234567890abcdef0123456789abcdef
  public_url: https://kasseapparat.example.com

payment_methods:
  - code: CASH
    name: Cash
  - code: SUMUP
    name: SumUp
```

The `payment_methods` list decides which buttons the POS offers. It is a YAML list, not a comma-separated string.

### Merchant Code (Merchant ID)

You can find your Merchant Code on [SumUp Settings](https://me.sumup.com/en-en/settings/developer), displayed directly below your company name.

### API Key

Log in to [SumUp Developer API Keys](https://me.sumup.com/en-en/settings/api-keys) and generate a new key.

Give it a descriptive name and store the key starting with `sup_sk` as `sumup.api_key`.

### Application ID

Log in to [Affiliate Keys](https://me.sumup.com/en-en/settings/affiliate-keys) and add a new Application ID (in the format `com.example.app`).

You can find more details in the [Getting Started guide](https://developer.sumup.com/terminal-payments/introduction/getting-started) of the SumUp SDK/API.

Store this ID as `sumup.application_id`.

### Affiliate Key

When creating the Application ID, you should copy the corresponding Affiliate Key from [Affiliate Keys](https://me.sumup.com/en-en/settings/affiliate-keys).

Store it as `sumup.affiliate_key`. The key starts with `sup_afk`.

### Public URL

If your Kasseapparat setup is publicly accessible via the internet, you can enter its public URL here. Otherwise, leave this field empty.

This URL is used for a **webhook**, which helps improve response times from the terminals.

## Pairing a Reader

Follow the instructions at [Pairing a Solo Reader](https://developer.sumup.com/terminal-payments/cloud-api#generate-pairing-code) until you see the **pairing code** displayed on the device.

Then go to the **Kasseapparat Admin** and navigate to **SumUp > Readers**.  
Click **"+ Pair"**, enter the pairing code, and provide a descriptive name for the reader.  
A meaningful name is especially useful if you are using multiple readers and need to distinguish between them.

After clicking **"Pair"**, the reader should now be successfully paired.

## Selecting a Reader

Once one or more readers are paired, you need to assign one to each terminal.

In the **SumUp > Readers** section of the admin interface, click **"Use this reader"** next to the reader you want to use.  
The selected reader will be marked as **"Selected"**.

You can now start accepting purchases in the frontend using the selected reader.

You can also choose the reader directly in the POS, without the admin:

- Open the user menu (top right) and tap **"Card reader: …"**.
- Or, when card payments cannot work, tap the hint below the payment buttons. It only appears when there is a problem: no reader selected, the selected reader is no longer paired, or its status is not "paired" (e.g. "expired"). The card payment button is disabled in that case.

The dialog lists all paired readers with their status; tap one to use it on this device. Like the selection in the admin, it is stored in this browser only, so every till keeps its own reader.

## Unpairing a Reader

A reader stays paired to your account. You cannot unpair a reader from the device itself!

So: do not forget to unpair the device after the party using the **Kasseapparat Admin**.
