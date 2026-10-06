package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

/*
	Organization X Employee Management and Payroll System

	No database is used.
	All information is stored in memory.

	The program demonstrates:
	1. Employee management
	2. Salary management
	3. Payment processing
	4. Reporting
	5. Black-box testing
	6. White-box testing
	7. Fault injection
*/

// ------------------------------------------------------------
// DATA MODEL
// ------------------------------------------------------------

type Employee struct {
	ID         string
	Name       string
	Department string
	Salary     float64
	Active     bool
}

// PayrollSystem stores employees and payments in memory.
type PayrollSystem struct {
	employees map[string]Employee
	payments  map[string]float64
}

// NewPayrollSystem creates an empty payroll system.
func NewPayrollSystem() *PayrollSystem {
	return &PayrollSystem{
		employees: make(map[string]Employee),
		payments:  make(map[string]float64),
	}
}

// ------------------------------------------------------------
// VALIDATION
// ------------------------------------------------------------

// validateEmployee is intentionally kept as a separate function.
// This makes it useful for white-box testing.
func validateEmployee(e Employee) error {
	if strings.TrimSpace(e.ID) == "" {
		return errors.New("employee ID cannot be empty")
	}

	if strings.TrimSpace(e.Name) == "" {
		return errors.New("employee name cannot be empty")
	}
	if !containsLetter(e.Name) {
		return errors.New("employee name must contain at least one letter")
	}

	if e.Salary <= 0 || math.IsNaN(e.Salary) || math.IsInf(e.Salary, 0) {
		return errors.New("salary must be greater than zero")
	}

	return nil
}

func containsLetter(value string) bool {
	for _, character := range value {
		if unicode.IsLetter(character) {
			return true
		}
	}
	return false
}

// ------------------------------------------------------------
// EMPLOYEE MANAGEMENT
// ------------------------------------------------------------

// AddEmployee registers a new employee.
func (ps *PayrollSystem) AddEmployee(e Employee) error {

	if err := validateEmployee(e); err != nil {
		return err
	}

	if _, exists := ps.employees[e.ID]; exists {
		return errors.New("employee already exists")
	}

	ps.employees[e.ID] = e

	return nil
}

// GetEmployee retrieves an employee by ID.
func (ps *PayrollSystem) GetEmployee(id string) (Employee, error) {

	employee, exists := ps.employees[id]

	if !exists {
		return Employee{}, errors.New("employee not found")
	}

	return employee, nil
}

// UpdateSalary changes an employee's salary.
func (ps *PayrollSystem) UpdateSalary(id string, salary float64) error {

	if salary <= 0 || math.IsNaN(salary) || math.IsInf(salary, 0) {
		return errors.New("salary must be greater than zero")
	}

	employee, exists := ps.employees[id]

	if !exists {
		return errors.New("employee not found")
	}

	employee.Salary = salary
	ps.employees[id] = employee

	return nil
}

// SetEmployeeActive activates or deactivates an employee.
func (ps *PayrollSystem) SetEmployeeActive(id string, active bool) error {

	employee, exists := ps.employees[id]

	if !exists {
		return errors.New("employee not found")
	}

	employee.Active = active
	ps.employees[id] = employee

	return nil
}

// ------------------------------------------------------------
// PAYMENT PROCESSING
// ------------------------------------------------------------

// ProcessPayment pays an active employee.
//
// Business rules:
// - Employee must exist.
// - Employee must be active.
// - Employee cannot be paid twice in one program execution.
func (ps *PayrollSystem) ProcessPayment(id string) (float64, error) {

	employee, exists := ps.employees[id]

	if !exists {
		return 0, errors.New("employee not found")
	}

	if !employee.Active {
		return 0, errors.New("employee is inactive")
	}

	if _, alreadyPaid := ps.payments[id]; alreadyPaid {
		return 0, errors.New("payment already processed")
	}

	ps.payments[id] = employee.Salary

	return employee.Salary, nil
}

// ------------------------------------------------------------
// REPORTING
// ------------------------------------------------------------

// GenerateReport produces a simple payroll report.
func (ps *PayrollSystem) GenerateReport() string {

	ids := make([]string, 0, len(ps.employees))

	for id := range ps.employees {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	activeCount := 0
	inactiveCount := 0
	activePayroll := 0.0

	for _, id := range ids {
		employee := ps.employees[id]

		if employee.Active {
			activeCount++
			activePayroll += employee.Salary
		} else {
			inactiveCount++
		}
	}

	paymentTotal := 0.0

	for _, amount := range ps.payments {
		paymentTotal += amount
	}

	var report strings.Builder

	fmt.Fprintf(&report, "========================================\n")
	fmt.Fprintf(&report, "ORGANIZATION X PAYROLL REPORT\n")
	fmt.Fprintf(&report, "========================================\n")
	fmt.Fprintf(&report, "Total employees: %d\n", len(ps.employees))
	fmt.Fprintf(&report, "Active employees: %d\n", activeCount)
	fmt.Fprintf(&report, "Inactive employees: %d\n", inactiveCount)
	fmt.Fprintf(&report, "Active monthly payroll: %.2f\n", activePayroll)
	fmt.Fprintf(&report, "Payments processed: %d\n", len(ps.payments))
	fmt.Fprintf(&report, "Payments total: %.2f\n", paymentTotal)
	fmt.Fprintf(&report, "----------------------------------------\n")

	for _, id := range ids {
		employee := ps.employees[id]

		status := "INACTIVE"

		if employee.Active {
			status = "ACTIVE"
		}

		fmt.Fprintf(
			&report,
			"%s | %s | %s | %.2f | %s\n",
			employee.ID,
			employee.Name,
			employee.Department,
			employee.Salary,
			status,
		)
	}

	fmt.Fprintf(&report, "========================================\n")

	return report.String()
}

// ------------------------------------------------------------
// TESTING UTILITIES
// ------------------------------------------------------------

type TestCase struct {
	ID          string
	Category    string
	Description string
	Run         func() (bool, string)
}

var totalTests int
var passedTests int
var failedTests int

func executeTest(tc TestCase, output io.Writer) {

	totalTests++

	fmt.Fprintf(
		output,
		"\n[%s] %s - %s\n",
		tc.Category,
		tc.ID,
		tc.Description,
	)

	passed, message := tc.Run()

	if passed {
		passedTests++
		fmt.Fprintln(output, "RESULT: PASS")
		fmt.Fprintf(output, "DETAIL: %s\n", message)
	} else {
		failedTests++
		fmt.Fprintln(output, "RESULT: FAIL")
		fmt.Fprintf(output, "DETAIL: %s\n", message)
	}
}

// ------------------------------------------------------------
// BLACK-BOX TESTS
// ------------------------------------------------------------

// BB-01: Add a valid employee.
func blackBox01() (bool, string) {

	ps := NewPayrollSystem()

	employee := Employee{
		ID:         "E001",
		Name:       "Alice",
		Department: "Finance",
		Salary:     3000,
		Active:     true,
	}

	err := ps.AddEmployee(employee)

	if err != nil {
		return false, fmt.Sprintf("Unexpected error: %v", err)
	}

	got, err := ps.GetEmployee("E001")

	if err != nil {
		return false, fmt.Sprintf("Employee could not be retrieved: %v", err)
	}

	if got.Name != "Alice" || got.Salary != 3000 {
		return false, fmt.Sprintf("Unexpected employee data: %+v", got)
	}

	return true, "Valid employee was successfully registered and retrieved."
}

// BB-02: Reject duplicate employee.
func blackBox02() (bool, string) {

	ps := NewPayrollSystem()

	employee := Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
		Active: true,
	}

	_ = ps.AddEmployee(employee)

	err := ps.AddEmployee(employee)

	if err == nil {
		return false, "Duplicate employee was accepted."
	}

	return true, fmt.Sprintf("Duplicate employee correctly rejected: %v", err)
}

// BB-03: Reject invalid employee data.
func blackBox03() (bool, string) {

	ps := NewPayrollSystem()

	employee := Employee{
		ID:     "",
		Name:   "Invalid Employee",
		Salary: 0,
		Active: true,
	}

	err := ps.AddEmployee(employee)

	if err == nil {
		return false, "Invalid employee was accepted."
	}

	return true, fmt.Sprintf("Invalid employee correctly rejected: %v", err)
}

// BB-04: Retrieve nonexistent employee.
func blackBox04() (bool, string) {

	ps := NewPayrollSystem()

	_, err := ps.GetEmployee("E999")

	if err == nil {
		return false, "Nonexistent employee was returned without an error."
	}

	return true, fmt.Sprintf("Nonexistent employee correctly rejected: %v", err)
}

// BB-05: Update employee salary.
func blackBox05() (bool, string) {

	ps := NewPayrollSystem()

	employee := Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
		Active: true,
	}

	_ = ps.AddEmployee(employee)

	err := ps.UpdateSalary("E001", 3500)

	if err != nil {
		return false, fmt.Sprintf("Salary update failed: %v", err)
	}

	updated, _ := ps.GetEmployee("E001")

	if updated.Salary != 3500 {
		return false, fmt.Sprintf(
			"Expected salary 3500.00 but got %.2f",
			updated.Salary,
		)
	}

	return true, "Employee salary successfully changed from 3000.00 to 3500.00."
}

// BB-06: Reject invalid salary update.
func blackBox06() (bool, string) {

	ps := NewPayrollSystem()

	employee := Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
		Active: true,
	}

	_ = ps.AddEmployee(employee)

	err := ps.UpdateSalary("E001", -500)

	if err == nil {
		return false, "Negative salary was accepted."
	}

	return true, fmt.Sprintf("Invalid salary correctly rejected: %v", err)
}

// BB-07: Process payment for active employee.
func blackBox07() (bool, string) {

	ps := NewPayrollSystem()

	employee := Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
		Active: true,
	}

	_ = ps.AddEmployee(employee)

	amount, err := ps.ProcessPayment("E001")

	if err != nil {
		return false, fmt.Sprintf("Payment failed: %v", err)
	}

	if amount != 3000 {
		return false, fmt.Sprintf(
			"Expected payment 3000.00 but got %.2f",
			amount,
		)
	}

	return true, "Active employee was paid exactly the current salary."
}

// BB-08: Reject payment for inactive employee.
func blackBox08() (bool, string) {

	ps := NewPayrollSystem()

	employee := Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
		Active: false,
	}

	_ = ps.AddEmployee(employee)

	_, err := ps.ProcessPayment("E001")

	if err == nil {
		return false, "Inactive employee was incorrectly paid."
	}

	return true, fmt.Sprintf("Inactive employee correctly rejected: %v", err)
}

// BB-09: Prevent duplicate payment.
func blackBox09() (bool, string) {

	ps := NewPayrollSystem()

	employee := Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
		Active: true,
	}

	_ = ps.AddEmployee(employee)

	_, firstErr := ps.ProcessPayment("E001")

	if firstErr != nil {
		return false, fmt.Sprintf("First payment failed: %v", firstErr)
	}

	_, secondErr := ps.ProcessPayment("E001")

	if secondErr == nil {
		return false, "Employee was paid twice."
	}

	return true, fmt.Sprintf("Duplicate payment correctly rejected: %v", secondErr)
}

// BB-10: Generate payroll report.
func blackBox10() (bool, string) {

	ps := NewPayrollSystem()

	_ = ps.AddEmployee(Employee{
		ID:         "E001",
		Name:       "Alice",
		Department: "Finance",
		Salary:     3000,
		Active:     true,
	})

	_ = ps.AddEmployee(Employee{
		ID:         "E002",
		Name:       "Bob",
		Department: "IT",
		Salary:     4000,
		Active:     false,
	})

	_, _ = ps.ProcessPayment("E001")

	report := ps.GenerateReport()

	expectedValues := []string{
		"Total employees: 2",
		"Active employees: 1",
		"Inactive employees: 1",
		"Active monthly payroll: 3000.00",
		"Payments processed: 1",
		"Payments total: 3000.00",
	}

	for _, expected := range expectedValues {
		if !strings.Contains(report, expected) {
			return false, fmt.Sprintf(
				"Report does not contain expected value: %s",
				expected,
			)
		}
	}

	return true, "Payroll report contains the expected employee and payment information."
}

// ------------------------------------------------------------
// WHITE-BOX TESTS
// ------------------------------------------------------------

// WB-01: validateEmployee - empty ID branch.
func whiteBox01() (bool, string) {

	e := Employee{
		ID:     "",
		Name:   "Alice",
		Salary: 3000,
	}

	err := validateEmployee(e)

	if err == nil {
		return false, "Expected empty-ID validation branch to return an error."
	}

	return true, "Empty employee ID branch was executed correctly."
}

// WB-02: validateEmployee - empty name branch.
func whiteBox02() (bool, string) {

	e := Employee{
		ID:     "E001",
		Name:   "",
		Salary: 3000,
	}

	err := validateEmployee(e)

	if err == nil {
		return false, "Expected empty-name validation branch to return an error."
	}

	return true, "Empty employee name branch was executed correctly."
}

// WB-03: validateEmployee - invalid salary branch.
func whiteBox03() (bool, string) {

	e := Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: -1,
	}

	err := validateEmployee(e)

	if err == nil {
		return false, "Expected invalid-salary branch to return an error."
	}

	return true, "Invalid salary validation branch was executed correctly."
}

// WB-04: validateEmployee - valid path.
func whiteBox04() (bool, string) {

	e := Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
	}

	err := validateEmployee(e)

	if err != nil {
		return false, fmt.Sprintf("Valid validation path returned an error: %v", err)
	}

	return true, "All validation conditions were satisfied and the valid path was executed."
}

// WB-05: AddEmployee duplicate branch.
func whiteBox05() (bool, string) {

	ps := NewPayrollSystem()

	e := Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
		Active: true,
	}

	_ = ps.AddEmployee(e)

	err := ps.AddEmployee(e)

	if err == nil {
		return false, "Duplicate branch was not triggered."
	}

	if err.Error() != "employee already exists" {
		return false, fmt.Sprintf(
			"Unexpected duplicate error: %v",
			err,
		)
	}

	return true, "AddEmployee duplicate branch returned the expected error."
}

// WB-06: UpdateSalary nonexistent employee branch.
func whiteBox06() (bool, string) {

	ps := NewPayrollSystem()

	err := ps.UpdateSalary("E999", 5000)

	if err == nil {
		return false, "Expected nonexistent-employee branch to return an error."
	}

	return true, "UpdateSalary nonexistent-employee branch executed correctly."
}

// WB-07: ProcessPayment nonexistent employee branch.
func whiteBox07() (bool, string) {

	ps := NewPayrollSystem()

	_, err := ps.ProcessPayment("E999")

	if err == nil {
		return false, "Expected ProcessPayment nonexistent branch to return an error."
	}

	return true, "ProcessPayment nonexistent-employee branch executed correctly."
}

// WB-08: ProcessPayment inactive employee branch.
func whiteBox08() (bool, string) {

	ps := NewPayrollSystem()

	_ = ps.AddEmployee(Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
		Active: false,
	})

	_, err := ps.ProcessPayment("E001")

	if err == nil {
		return false, "Expected inactive branch to return an error."
	}

	return true, "ProcessPayment inactive-employee branch executed correctly."
}

// WB-09: ProcessPayment already-paid branch.
func whiteBox09() (bool, string) {

	ps := NewPayrollSystem()

	_ = ps.AddEmployee(Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
		Active: true,
	})

	_, _ = ps.ProcessPayment("E001")

	_, err := ps.ProcessPayment("E001")

	if err == nil {
		return false, "Expected already-paid branch to return an error."
	}

	return true, "ProcessPayment already-paid branch executed correctly."
}

// WB-10: ProcessPayment successful branch.
func whiteBox10() (bool, string) {

	ps := NewPayrollSystem()

	_ = ps.AddEmployee(Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
		Active: true,
	})

	amount, err := ps.ProcessPayment("E001")

	if err != nil {
		return false, fmt.Sprintf("Successful branch returned an error: %v", err)
	}

	if amount != 3000 {
		return false, fmt.Sprintf(
			"Successful branch returned %.2f instead of 3000.00",
			amount,
		)
	}

	return true, "Successful payment branch returned the expected amount."
}

// ------------------------------------------------------------
// FAULT INJECTION
// ------------------------------------------------------------

// faultyCalculatePayment intentionally introduces a defect.
// It adds 100 to the actual salary.
func faultyCalculatePayment(e Employee) float64 {
	return e.Salary + 100
}

// This test proves that our test logic can detect an injected fault.
func runFaultInjectionTest(output io.Writer) {

	fmt.Fprintln(output, "\n========================================")
	fmt.Fprintln(output, "FAULT INJECTION TEST")
	fmt.Fprintln(output, "========================================")

	employee := Employee{
		ID:     "E001",
		Name:   "Alice",
		Salary: 3000,
		Active: true,
	}

	expected := employee.Salary
	actual := faultyCalculatePayment(employee)

	fmt.Fprintf(output, "Expected payment: %.2f\n", expected)
	fmt.Fprintf(output, "Faulty payment:   %.2f\n", actual)

	if actual != expected {
		fmt.Fprintln(output, "RESULT: PASS")
		fmt.Fprintln(output,
			"DETAIL: The test successfully detected the injected payroll defect.",
		)
	} else {
		fmt.Fprintln(output, "RESULT: FAIL")
		fmt.Fprintln(output, "DETAIL: The injected defect was not detected.")
	}
}

// ------------------------------------------------------------
// DEMONSTRATION OF THE APPLICATION
// ------------------------------------------------------------

func demonstrateApplication() {

	fmt.Println("\n========================================")
	fmt.Println("APPLICATION DEMONSTRATION")
	fmt.Println("========================================")

	ps := NewPayrollSystem()

	fmt.Println("\nAdding employees...")

	err := ps.AddEmployee(Employee{
		ID:         "E001",
		Name:       "Alice",
		Department: "Finance",
		Salary:     3000,
		Active:     true,
	})

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Println("SUCCESS: Alice added.")
	}

	err = ps.AddEmployee(Employee{
		ID:         "E002",
		Name:       "Bob",
		Department: "IT",
		Salary:     4000,
		Active:     true,
	})

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Println("SUCCESS: Bob added.")
	}

	fmt.Println("\nUpdating Alice's salary...")

	err = ps.UpdateSalary("E001", 3500)

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Println("SUCCESS: Alice's salary updated to 3500.00.")
	}

	fmt.Println("\nProcessing payments...")

	amount, err := ps.ProcessPayment("E001")

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("SUCCESS: Alice paid %.2f\n", amount)
	}

	amount, err = ps.ProcessPayment("E002")

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("SUCCESS: Bob paid %.2f\n", amount)
	}

	fmt.Println("\nGenerating report...")
	fmt.Println(ps.GenerateReport())
}

func readPrompt(scanner *bufio.Scanner, output io.Writer, prompt string) (string, bool) {
	fmt.Fprint(output, prompt)
	if !scanner.Scan() {
		return "", false
	}
	return strings.TrimSpace(scanner.Text()), true
}

func readValidatedPrompt(scanner *bufio.Scanner, output io.Writer, prompt string, validate func(string) error) (string, bool) {
	for {
		value, ok := readPrompt(scanner, output, prompt)
		if !ok {
			return "", false
		}
		if err := validate(value); err != nil {
			fmt.Fprintf(output, "Invalid input: %v\n", err)
			continue
		}
		return value, true
	}
}

func readSalaryPrompt(scanner *bufio.Scanner, output io.Writer, prompt string) (float64, bool) {
	for {
		value, ok := readPrompt(scanner, output, prompt)
		if !ok {
			return 0, false
		}
		salary, err := strconv.ParseFloat(value, 64)
		if err != nil || salary <= 0 || math.IsNaN(salary) || math.IsInf(salary, 0) {
			fmt.Fprintln(output, "Invalid input: salary must be a valid number greater than zero.")
			continue
		}
		return salary, true
	}
}

func readExistingEmployeeID(scanner *bufio.Scanner, output io.Writer, ps *PayrollSystem) (string, bool) {
	return readValidatedPrompt(scanner, output, "Employee ID: ", func(id string) error {
		if id == "" {
			return errors.New("employee ID cannot be empty")
		}
		if _, exists := ps.employees[id]; !exists {
			return errors.New("employee not found")
		}
		return nil
	})
}

func runInteractiveApplication(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	ps := NewPayrollSystem()

	fmt.Fprintln(output, "Organization X Employee Management & Payroll")
	for {
		fmt.Fprintln(output, "\n1. Add employee")
		fmt.Fprintln(output, "2. Update salary")
		fmt.Fprintln(output, "3. Change employee active status")
		fmt.Fprintln(output, "4. Process payment")
		fmt.Fprintln(output, "5. Show payroll report")
		fmt.Fprintln(output, "6. Run automated tests")
		fmt.Fprintln(output, "0. Exit")

		choice, ok := readPrompt(scanner, output, "Choose an option: ")
		if !ok {
			return
		}

		switch choice {
		case "1":
			id, ok := readValidatedPrompt(scanner, output, "Employee ID: ", func(id string) error {
				if id == "" {
					return errors.New("employee ID cannot be empty")
				}
				if _, exists := ps.employees[id]; exists {
					return errors.New("employee already exists")
				}
				return nil
			})
			if !ok {
				return
			}
			name, ok := readValidatedPrompt(scanner, output, "Employee name: ", func(name string) error {
				if name == "" {
					return errors.New("employee name cannot be empty")
				}
				if !containsLetter(name) {
					return errors.New("employee name must contain at least one letter")
				}
				return nil
			})
			if !ok {
				return
			}
			department, ok := readValidatedPrompt(scanner, output, "Department: ", func(department string) error {
				if department == "" {
					return errors.New("department cannot be empty")
				}
				if !containsLetter(department) {
					return errors.New("department must contain at least one letter")
				}
				return nil
			})
			if !ok {
				return
			}
			salary, ok := readSalaryPrompt(scanner, output, "Salary: ")
			if !ok {
				return
			}

			err := ps.AddEmployee(Employee{
				ID:         id,
				Name:       name,
				Department: department,
				Salary:     salary,
				Active:     true,
			})
			if err != nil {
				fmt.Fprintf(output, "ERROR: %v\n", err)
			} else {
				fmt.Fprintf(output, "SUCCESS: %s added.\n", id)
			}
		case "2":
			id, ok := readExistingEmployeeID(scanner, output, ps)
			if !ok {
				return
			}
			salary, ok := readSalaryPrompt(scanner, output, "New salary: ")
			if !ok {
				return
			}
			if err := ps.UpdateSalary(id, salary); err != nil {
				fmt.Fprintf(output, "ERROR: %v\n", err)
			} else {
				fmt.Fprintf(output, "SUCCESS: salary updated for %s.\n", id)
			}
		case "3":
			id, ok := readExistingEmployeeID(scanner, output, ps)
			if !ok {
				return
			}
			status, ok := readValidatedPrompt(scanner, output, "Set active? (y/n): ", func(status string) error {
				switch strings.ToLower(status) {
				case "y", "yes", "n", "no":
					return nil
				default:
					return errors.New("enter y or n")
				}
			})
			if !ok {
				return
			}
			active := strings.EqualFold(status, "y") || strings.EqualFold(status, "yes")
			if err := ps.SetEmployeeActive(id, active); err != nil {
				fmt.Fprintf(output, "ERROR: %v\n", err)
			} else {
				fmt.Fprintf(output, "SUCCESS: active status updated for %s.\n", id)
			}
		case "4":
			id, ok := readExistingEmployeeID(scanner, output, ps)
			if !ok {
				return
			}
			amount, err := ps.ProcessPayment(id)
			if err != nil {
				fmt.Fprintf(output, "ERROR: %v\n", err)
			} else {
				fmt.Fprintf(output, "SUCCESS: payment of %.2f processed for %s.\n", amount, id)
			}
		case "5":
			fmt.Fprintln(output, ps.GenerateReport())
		case "6":
			runAutomatedTests(output)
		case "0":
			fmt.Fprintln(output, "Goodbye.")
			return
		default:
			fmt.Fprintln(output, "Invalid option. Choose 0 through 6.")
		}
	}
}

// ------------------------------------------------------------
// MAIN
// ------------------------------------------------------------

func runAutomatedTests(output io.Writer) {
	totalTests = 0
	passedTests = 0
	failedTests = 0

	fmt.Fprintln(output, "========================================")
	fmt.Fprintln(output, "ORGANIZATION X AUTOMATED TEST EXECUTION")
	fmt.Fprintln(output, "========================================")
	fmt.Fprintln(output, "\n\n========================================")
	fmt.Fprintln(output, "BLACK-BOX TESTS")
	fmt.Fprintln(output, "========================================")

	blackBoxTests := []TestCase{
		{
			ID:          "BB-01",
			Category:    "BLACK-BOX",
			Description: "Register a valid employee",
			Run:         blackBox01,
		},
		{
			ID:          "BB-02",
			Category:    "BLACK-BOX",
			Description: "Reject duplicate employee ID",
			Run:         blackBox02,
		},
		{
			ID:          "BB-03",
			Category:    "BLACK-BOX",
			Description: "Reject invalid employee information",
			Run:         blackBox03,
		},
		{
			ID:          "BB-04",
			Category:    "BLACK-BOX",
			Description: "Retrieve nonexistent employee",
			Run:         blackBox04,
		},
		{
			ID:          "BB-05",
			Category:    "BLACK-BOX",
			Description: "Update employee salary",
			Run:         blackBox05,
		},
		{
			ID:          "BB-06",
			Category:    "BLACK-BOX",
			Description: "Reject invalid salary update",
			Run:         blackBox06,
		},
		{
			ID:          "BB-07",
			Category:    "BLACK-BOX",
			Description: "Process payment for active employee",
			Run:         blackBox07,
		},
		{
			ID:          "BB-08",
			Category:    "BLACK-BOX",
			Description: "Reject payment for inactive employee",
			Run:         blackBox08,
		},
		{
			ID:          "BB-09",
			Category:    "BLACK-BOX",
			Description: "Prevent duplicate payment",
			Run:         blackBox09,
		},
		{
			ID:          "BB-10",
			Category:    "BLACK-BOX",
			Description: "Generate payroll report",
			Run:         blackBox10,
		},
	}

	for _, test := range blackBoxTests {
		executeTest(test, output)
	}

	fmt.Fprintln(output, "\n\n========================================")
	fmt.Fprintln(output, "WHITE-BOX TESTS")
	fmt.Fprintln(output, "========================================")

	whiteBoxTests := []TestCase{
		{
			ID:          "WB-01",
			Category:    "WHITE-BOX",
			Description: "Execute empty employee-ID validation branch",
			Run:         whiteBox01,
		},
		{
			ID:          "WB-02",
			Category:    "WHITE-BOX",
			Description: "Execute empty employee-name validation branch",
			Run:         whiteBox02,
		},
		{
			ID:          "WB-03",
			Category:    "WHITE-BOX",
			Description: "Execute invalid-salary validation branch",
			Run:         whiteBox03,
		},
		{
			ID:          "WB-04",
			Category:    "WHITE-BOX",
			Description: "Execute valid employee validation path",
			Run:         whiteBox04,
		},
		{
			ID:          "WB-05",
			Category:    "WHITE-BOX",
			Description: "Execute duplicate employee branch",
			Run:         whiteBox05,
		},
		{
			ID:          "WB-06",
			Category:    "WHITE-BOX",
			Description: "Execute nonexistent employee salary-update branch",
			Run:         whiteBox06,
		},
		{
			ID:          "WB-07",
			Category:    "WHITE-BOX",
			Description: "Execute nonexistent employee payment branch",
			Run:         whiteBox07,
		},
		{
			ID:          "WB-08",
			Category:    "WHITE-BOX",
			Description: "Execute inactive employee payment branch",
			Run:         whiteBox08,
		},
		{
			ID:          "WB-09",
			Category:    "WHITE-BOX",
			Description: "Execute duplicate payment branch",
			Run:         whiteBox09,
		},
		{
			ID:          "WB-10",
			Category:    "WHITE-BOX",
			Description: "Execute successful payment branch",
			Run:         whiteBox10,
		},
	}

	for _, test := range whiteBoxTests {
		executeTest(test, output)
	}

	runFaultInjectionTest(output)

	fmt.Fprintln(output, "\n========================================")
	fmt.Fprintln(output, "FINAL TEST SUMMARY")
	fmt.Fprintln(output, "========================================")

	fmt.Fprintf(output, "Total automated tests: %d\n", totalTests)
	fmt.Fprintf(output, "Passed:                %d\n", passedTests)
	fmt.Fprintf(output, "Failed:                %d\n", failedTests)

	if failedTests == 0 {
		fmt.Fprintln(output, "OVERALL RESULT: PASS")
	} else {
		fmt.Fprintln(output, "OVERALL RESULT: FAIL")
	}

	fmt.Fprintln(output, "========================================")
}

func main() {
	runInteractiveApplication(os.Stdin, os.Stdout)
}
