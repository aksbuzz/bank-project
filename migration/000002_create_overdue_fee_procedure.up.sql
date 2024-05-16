CREATE OR REPLACE PROCEDURE calculate_overdue_fees()
LANGUAGE plpgsql
AS $$
BEGIN
  UPDATE loans
  SET overdue_fee = (CURRENT_DATE - due_date) * 0.2 -- $0.2 per day overdue
  WHERE return_date IS NULL AND CURRENT_DATE > due_date;
END;
$$
