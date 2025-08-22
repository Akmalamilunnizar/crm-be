import { string } from "yup";

export const defaultStringValidation = (noOfChars = 3, isEmail = false) => {
  const base = string()
    .min(noOfChars, `Must be at least ${noOfChars} characters`)
    .required("Required");

  if (isEmail) {
    return base.email("Invalid email");
  }

  return base;
};
