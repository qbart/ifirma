# ifirma

Simple ifirma CLI tool. Invoices are configured to use 23% VAT.

## config

Create `~/.config/ifirma.hcl` file with the following content:

```hcl
payment {
  # spaces will be removed 
  bank = "63 xxxx xxxx xxxx xxxx xxxx xxxx"
}

invoice "invoice_ID" {
  to = "COMPANY" # contractor ID

  issued_at = "2020-11-03" # if empty it will take the last day of prev month
  sold_at   = "2020-11-03" # if empty it will take the last day of prev month


  # only one pos is supported at the moment
  pos {
    name     = "Name of the service/product"
    quantity = 1
    unit     = "szt"
    gtu      = "12" # one of ["01",...,"12"] or "BRAK" - this parameter is optional
  }
}

invoice "invoice_ID2" {
  to = "OTHER-COMPANY" # contractor ID

  ...
}

...

```

## usage

CLI requires following envs to be defined:
- `IFIRMA_FV_TOKEN`
- `IFIRMA_EMAIL`

```
ifirma -invoice INVOICE_ID -net_price 5000.00
ifirma -invoice INVOICE_ID -net_price 5000.00 -config ~/diffrent_path_to_config/ifirma.hcl
```

## contribution guide

Just ask before adding new code.
