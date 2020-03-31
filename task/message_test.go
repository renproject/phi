package task_test

import (
	"context"
	"testing/quick"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/phi/task"
)

var _ = Describe("Messages", func() {
	Context("when building a message", func() {
		It("should return the expected request", func() {
			f := func(expectedRequest int) bool {
				message := task.NewMessage(expectedRequest)
				request := message.Request()
				Expect(request).To(Equal(expectedRequest))
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})

	Context("when waiting for a response", func() {
		It("should return the expected response", func() {
			f := func(expectedResponse int) bool {
				message := task.NewMessage(nil)
				message.Respond(expectedResponse)
				response, err := message.Wait(context.Background())
				Expect(err).ToNot(HaveOccurred())
				Expect(response).To(Equal(expectedResponse))
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})
})
