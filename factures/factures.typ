#set page(
  margin: (top: 1.5cm, bottom: 1.5cm, left: 2cm, right: 2cm),
  footer: [
    #set text(size: 8pt, fill: gray)
    #align(center)[
      superpdp.fr 
    ]
  ]
)

#set text(
  font: "Liberation Sans",
  size: 10pt,
  fill: rgb("#2b2b2b")
)

#let accent = rgb("#4f6bed")
#let light-gray = rgb("#f8f9fa")
#let border-color = rgb("#e5e7eb")

#let data = json("invoice.json")

#grid(
  columns: (1fr, 1fr),
  [
    #set align(left)
    #box(
      fill: rgb("#4244d6"),
      radius: 6pt,
      inset: 10pt
    )[
      #image("superpdpd.png", width: 3cm)
    ]

    #text(size: 14pt, weight: "bold")[#data.seller.name] \
    #text(size: 9pt, fill: gray)[
      #data.seller.legal_registration_identifier.value (#data.seller.legal_registration_identifier.scheme) \
    ]
  ],
  [
    #set align(right)
    #text(size: 24pt, weight: "bold", fill: accent)[FACTURE] \
    #text(size: 11pt, weight: "medium")[#data.number] \
    #text(size: 9pt, fill: gray)[Issue date : #data.issue_date]
  ]
)

#v(1.5cm)

#grid(
  columns: (1fr, 1fr),
  gutter: 2cm,
  [
    #set align(left)
    #text(fill: accent, weight: "bold", size: 9pt)[SELLER] \
    #line(length: 100%, stroke: 0.5pt + border-color)
    #v(0.1cm)
    #text(weight: "bold")[#data.seller.name] \
    #data.seller.postal_address.country_code
  ],
  [
    #set align(left)
    #text(fill: accent, weight: "bold", size: 9pt)[BUYER] \
    #line(length: 100%, stroke: 0.5pt + border-color)
    #v(0.1cm)
    #text(weight: "bold")[#data.buyer.name] \
    #data.buyer.postal_address.country_code
  ]
)

#v(1cm)

#rect(
  fill: light-gray,
  inset: 12pt,
  radius: 4pt,
  width: 100%,
  grid(
    columns: (1fr, 1fr, 1fr),
    [ #text(gray, size: 8pt)[Due date] \ #data.payment_due_date ],
    [ #text(gray, size: 8pt)[Currency] \ #data.currency_code ],
    [ #text(gray, size: 8pt)[Payement method] \ Transfer ]
  )
)

#v(0.5cm)

#table(
  columns: (1fr, 80pt, 80pt, 80pt),
  stroke: (x, y) => if y == 0 { (bottom: 1pt + accent) } else { (bottom: 0.5pt + border-color) },
  inset: 10pt,
  fill: (x, y) => if y == 0 { none } else if calc.even(y) { white } else { white },
  
  table.header(
    text(weight: "bold", size: 9pt)[Description],
    text(weight: "bold", size: 9pt)[Unit price],
    text(weight: "bold", size: 9pt)[Quantity],
    text(weight: "bold", size: 9pt)[Total excluding tax],
  ),
  ..for produit in data.lines {
    (
      [#produit.item_information.name], 
      [#produit.price_details.item_net_price], 
      [#produit.invoiced_quantity], 
      [#produit.net_amount]
    )
  }
)

#v(0.5cm)
#grid(
  columns: (1fr, 180pt),
  [
    #if data.notes.len() > 0 [
      #text(size: 8pt, fill: gray)[
        #for note in data.notes {
          [#note.subject_code : #note.note \ ]
        }
      ]
    ]
  ],
  [
    #set align(right)
    #table(
      columns: (1fr, 1fr),
      stroke: none,
      inset: 5pt,
      [Total HT], [#data.totals.total_without_vat #data.currency_code],
      [TVA], [#data.totals.total_vat_amount.value #data.currency_code],
      line(length: 100%, stroke: 0.5pt + border-color), line(length: 100%, stroke: 0.5pt + border-color),
      text(weight: "bold")[Net à payer], text(weight: "bold", fill: accent, size: 12pt)[#data.totals.amount_due_for_payment #data.currency_code]
    )
  ]
)