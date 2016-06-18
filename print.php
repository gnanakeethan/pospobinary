<?php
/* ASCII constants */
const ESC = "\x1b";
const GS="\x1d";
const NUL="\x00";

/* Output an example receipt */
echo ESC."@"; // Reset to defaults
echo GS."k"."\x04"."123123".NUL; // Print barcode
echo ESC."d".chr(1); // Blank line
exit(0);
