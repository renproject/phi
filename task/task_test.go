package task_test

import (
	"context"
	"testing/quick"
	"time"

	"github.com/renproject/phi/task"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	Context("when sending messages", func() {
		Context("when waiting for a response", func() {
			It("should return the expected response", func() {
				f := func(expectedResponse int) bool {
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()

					t := task.New(task.HandlerFunc(func(message task.Message) { message.Respond(expectedResponse) }), 1)
					go t.Run(ctx)

					response := 0
					err := t.SendAndWait(ctx, nil, &response)
					Expect(err).ToNot(HaveOccurred())
					Expect(response).To(Equal(expectedResponse))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})

			Context("when the response is not a pointer", func() {
				It("should panic", func() {
					f := func(expectedResponse int) bool {
						ctx, cancel := context.WithCancel(context.Background())
						defer cancel()

						t := task.New(task.HandlerFunc(func(message task.Message) {
							message.Respond(expectedResponse)
						}), 1)
						go t.Run(ctx)

						Expect(func() { t.SendAndWait(ctx, nil, struct{}{}) }).To(Panic())
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when the context is done", func() {
				It("should return an error", func() {
					f := func(expectedResponse int) bool {
						ctx, cancel := context.WithCancel(context.Background())
						defer cancel()

						t := task.New(task.HandlerFunc(func(message task.Message) {
							time.Sleep(10 * time.Millisecond)
							message.Respond(expectedResponse)
						}), 0) // Capacity is set to zero so that context cancellation can kick in.
						go t.Run(ctx)

						sendCtx, sendCancel := context.WithTimeout(context.Background(), time.Millisecond)
						defer sendCancel()

						response := 0
						err := t.SendAndWait(sendCtx, nil, &response)
						Expect(err).To(HaveOccurred())
						Expect(response).To(Equal(0))
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when the response is nil", func() {
				It("should panic", func() {
					f := func(expectedResponse int) bool {
						ctx, cancel := context.WithCancel(context.Background())
						defer cancel()

						t := task.New(task.HandlerFunc(func(message task.Message) {
							time.Sleep(2 * time.Millisecond)
							message.Respond(expectedResponse)
						}), 0) // Capacity is set to zero so that context cancellation can kick in.
						go t.Run(ctx)

						sendCtx, sendCancel := context.WithTimeout(context.Background(), time.Millisecond)
						defer sendCancel()

						Expect(func() { t.SendAndWait(sendCtx, struct{}{}, nil) }).To(Panic())
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})
		})

		Context("when the context is done", func() {
			It("should return an error", func() {
				f := func() bool {
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()

					t := task.New(task.HandlerFunc(func(message task.Message) {
						message.Respond(struct{}{})
					}), 0) // Capacity is set to zero so that context cancellation can kick in.
					go t.Run(ctx)

					sendCtx, sendCancel := context.WithCancel(context.Background())
					sendCancel()

					err := t.SendAndWait(sendCtx, struct{}{}, nil)
					Expect(err).To(HaveOccurred())
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})
	})
})
