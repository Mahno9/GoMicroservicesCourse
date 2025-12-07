package e2e

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

var _ = Describe("Inventory Service", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc

		conn            *grpc.ClientConn
		invnetoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		// Create gRPC connection to the inventory service
		var err error
		conn, err = grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "🟥 Failed to create gRPC connection")

		invnetoryClient = inventoryV1.NewInventoryServiceClient(conn)

		// Clean up the parts collection before each test
		err = env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "🟥 Failed to clear parts collection")
	})

	AfterEach(func() {
		if conn != nil {
			err := conn.Close()
			Expect(err).ToNot(HaveOccurred(), "🟥 Failed to close gRPC connection")
		}
		cancel()
	})

	Describe("GetPart", func() {
		Context("when part exists", func() {
			It("should return the part successfully", func() {
				// Arrange: Insert a test part
				partUUID, err := env.InsertTestPart(ctx)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")
				Expect(partUUID).NotTo(BeEmpty(), "🟥 Test part UUID is empty")

				// Act: Get the part
				resp, err := invnetoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
					Uuid: partUUID,
				})

				// Assert
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to get part")
				Expect(resp).ToNot(BeNil(), "🟥 Response is nil")
				Expect(resp.Part).ToNot(BeNil(), "🟥 Part is nil")
				Expect(resp.Part.Uuid).To(Equal(partUUID), "🟥 Part UUID does not match")
				Expect(resp.Part.Name).NotTo(BeEmpty(), "🟥 Part name is empty")
				Expect(resp.Part.Description).NotTo(BeEmpty(), "🟥 Part description is empty")
				Expect(resp.Part.Price).To(BeNumerically(">", 0), "🟥 Part price is not positive")
				Expect(resp.Part.StockQuantity).To(BeNumerically(">", 0), "🟥 Part stock quantity is not positive")
				Expect(resp.Part.Category).To(Equal(inventoryV1.Category_ENGINE), "🟥 Part category is not ENGINE")
				Expect(resp.Part.Dimensions).NotTo(BeNil(), "🟥 Part dimensions is nil")
				Expect(resp.Part.Manufacturer).NotTo(BeNil(), "🟥 Part manufacturer is nil")
				Expect(resp.Part.Tags).NotTo(BeEmpty(), "🟥 Part tags is empty")
				Expect(resp.Part.CreatedAt).NotTo(BeNil(), "🟥 Part created at is nil")
			})

			It("should return part with correct metadata", func() {
				// Arrange: Insert a test part with specific data
				testInfo := env.GetTestPartInfo()
				partUUID, err := env.InsertTestPartWithData(ctx, testInfo)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")

				// Act: Get the part
				resp, err := invnetoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
					Uuid: partUUID,
				})

				// Assert
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to get part")
				Expect(resp).ToNot(BeNil(), "🟥 Response is nil")
				Expect(resp.Part).ToNot(BeNil(), "🟥 Part is nil")
				Expect(resp.Part.Name).To(Equal(testInfo.Name), "🟥 Part name does not match")
				Expect(resp.Part.Description).To(Equal(testInfo.Description), "🟥 Part description does not match")
				Expect(resp.Part.Price).To(Equal(testInfo.Price), "🟥 Part price does not match")
				Expect(resp.Part.StockQuantity).To(Equal(testInfo.StockQuantity), "🟥 Part stock quantity does not match")
				Expect(resp.Part.Category).To(Equal(testInfo.Category), "🟥 Part category does not match")
				Expect(resp.Part.Manufacturer.Name).To(Equal(testInfo.Manufacturer.Name), "🟥 Part manufacturer name does not match")
				Expect(resp.Part.Manufacturer.Country).To(Equal(testInfo.Manufacturer.Country), "🟥 Part manufacturer country does not match")
				Expect(resp.Part.Tags).To(Equal(testInfo.Tags), "🟥 Part tags do not match")
			})
		})

		Context("when part does not exist", func() {
			It("should return NotFound error", func() {
				// Arrange: Use a non-existent UUID
				nonExistentUUID := "00000000-0000-0000-0000-000000000000"

				// Act: Try to get a non-existent part
				resp, err := invnetoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
					Uuid: nonExistentUUID,
				})

				// Assert
				Expect(err).To(HaveOccurred(), "🟥 Failed to get part")
				Expect(resp).To(BeNil(), "🟥 Response is nil")

				st, ok := status.FromError(err)
				Expect(ok).To(BeTrue(), "🟥 Status is nil")
				Expect(st.Code()).To(Equal(codes.NotFound), "🟥 Status code is not NotFound")
			})
		})

		Context("when UUID is invalid", func() {
			It("should return InvalidArgument error for empty UUID", func() {
				// Act: Try to get a part with empty UUID
				resp, err := invnetoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
					Uuid: "",
				})

				// Assert
				Expect(err).To(HaveOccurred(), "🟥 Failed to get part")
				Expect(resp).To(BeNil(), "🟥 Response is nil")

				st, ok := status.FromError(err)
				Expect(ok).To(BeTrue(), "🟥 Status is nil")
				Expect(st.Code()).To(Equal(codes.InvalidArgument), "🟥 Status code is not InvalidArgument")
			})
		})
	})

	Describe("ListParts", func() {
		Context("when no filters are applied", func() {
			It("should return all parts", func() {
				// Arrange: Insert multiple test parts
				part1UUID, err := env.InsertTestPart(ctx)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")
				part2UUID, err := env.InsertTestPart(ctx)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")

				// Act: List all parts
				resp, err := invnetoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})

				// Assert
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to list parts")
				Expect(resp).ToNot(BeNil(), "🟥 Response is nil")
				Expect(resp.Parts).To(HaveLen(2), "🟥 Parts count is not 2")

				// Verify that both parts are in the response
				partUUIDs := []string{resp.Parts[0].Uuid, resp.Parts[1].Uuid}
				Expect(partUUIDs).To(ContainElement(part1UUID))
				Expect(partUUIDs).To(ContainElement(part2UUID))
			})

			It("should return empty list when no parts exist", func() {
				// Act: List parts when collection is empty
				resp, err := invnetoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})

				// Assert
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to list parts")
				Expect(resp).ToNot(BeNil(), "🟥 Response is nil")
				Expect(resp.Parts).To(BeEmpty(), "🟥 Parts count is not 0")
			})
		})

		Context("when filtering by UUIDs", func() {
			It("should return only parts with matching UUIDs", func() {
				// Arrange: Insert multiple parts
				part1UUID, err := env.InsertTestPart(ctx)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")
				part2UUID, err := env.InsertTestPart(ctx)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")
				_, err = env.InsertTestPart(ctx)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")

				// Act: List parts filtering by specific UUIDs
				resp, err := invnetoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
					Filter: &inventoryV1.PartsFilter{
						Uuids: []string{part1UUID, part2UUID},
					},
				})

				// Assert
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to list parts")
				Expect(resp).ToNot(BeNil(), "🟥 Response is nil")
				Expect(resp.Parts).To(HaveLen(2), "🟥 Parts count is not 2")

				partUUIDs := []string{resp.Parts[0].Uuid, resp.Parts[1].Uuid}
				Expect(partUUIDs).To(ContainElement(part1UUID))
				Expect(partUUIDs).To(ContainElement(part2UUID))
			})

			It("should return empty list when UUIDs don't match", func() {
				// Arrange: Insert a test part
				_, err := env.InsertTestPart(ctx)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")

				// Act: List parts filtering by non-existent UUID
				resp, err := invnetoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
					Filter: &inventoryV1.PartsFilter{
						Uuids: []string{"00000000-0000-0000-0000-000000000000"},
					},
				})

				// Assert
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to list parts")
				Expect(resp).ToNot(BeNil(), "🟥 Response is nil")
				Expect(resp.Parts).To(BeEmpty(), "🟥 Parts count is not 0")
			})
		})

		Context("when filtering by names", func() {
			It("should return only parts with matching names", func() {
				// Arrange: Insert parts with specific names
				testInfo1 := &inventoryV1.Part{
					Name:          "Boeing Wing",
					Description:   "Main wing component",
					Price:         5000,
					StockQuantity: 5,
					Category:      inventoryV1.Category_WING,
					Dimensions: &inventoryV1.Dimensions{
						Length: 10, Width: 5, Height: 2, Weight: 100,
					},
					Manufacturer: &inventoryV1.Manufacturer{
						Name: "Boeing", Country: "USA", Website: "https://boeing.com",
					},
					Tags: []string{"wing", "aircraft"},
				}
				partUUID1, err := env.InsertTestPartWithData(ctx, testInfo1)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")

				testInfo2 := &inventoryV1.Part{
					Name:          "Airbus Engine",
					Description:   "Jet engine",
					Price:         10000,
					StockQuantity: 3,
					Category:      inventoryV1.Category_ENGINE,
					Dimensions: &inventoryV1.Dimensions{
						Length: 5, Width: 3, Height: 3, Weight: 500,
					},
					Manufacturer: &inventoryV1.Manufacturer{
						Name: "Airbus", Country: "France", Website: "https://airbus.com",
					},
					Tags: []string{"engine", "aircraft"},
				}
				_, err = env.InsertTestPartWithData(ctx, testInfo2)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")

				// Act: Filter by name
				resp, err := invnetoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
					Filter: &inventoryV1.PartsFilter{
						Names: []string{"Boeing Wing"},
					},
				})

				// Assert
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to list parts")
				Expect(resp).ToNot(BeNil(), "🟥 Response is nil")
				Expect(resp.Parts).To(HaveLen(1), "🟥 Parts count is not 1")
				Expect(resp.Parts[0].Uuid).To(Equal(partUUID1), "🟥 Part UUID does not match")
				Expect(resp.Parts[0].Name).To(Equal("Boeing Wing"), "🟥 Part name does not match")
			})
		})

		Context("when filtering by categories", func() {
			It("should return only parts with matching categories", func() {
				// Arrange: Insert parts with different categories
				engineInfo := &inventoryV1.Part{
					Name:          "Engine Part",
					Description:   "Engine component",
					Price:         3000,
					StockQuantity: 10,
					Category:      inventoryV1.Category_ENGINE,
					Dimensions: &inventoryV1.Dimensions{
						Length: 5, Width: 5, Height: 5, Weight: 50,
					},
					Manufacturer: &inventoryV1.Manufacturer{
						Name: "EngCo", Country: "USA", Website: "https://engco.com",
					},
					Tags: []string{"engine"},
				}
				engineUUID, err := env.InsertTestPartWithData(ctx, engineInfo)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")

				wingInfo := env.GetTestPartInfo() // This uses WING category
				_, err = env.InsertTestPartWithData(ctx, wingInfo)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")

				// Act: Filter by ENGINE category
				resp, err := invnetoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
					Filter: &inventoryV1.PartsFilter{
						Categories: []inventoryV1.Category{inventoryV1.Category_ENGINE},
					},
				})

				// Assert
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to list parts")
				Expect(resp).ToNot(BeNil(), "🟥 Response is nil")
				Expect(resp.Parts).To(HaveLen(1), "🟥 Parts count is not 1")
				Expect(resp.Parts[0].Uuid).To(Equal(engineUUID), "🟥 Part UUID does not match")
				Expect(resp.Parts[0].Category).To(Equal(inventoryV1.Category_ENGINE), "🟥 Part category does not match")
			})
		})

		Context("when filtering by manufacturer country", func() {
			It("should return only parts from specific countries", func() {
				// Arrange: Insert parts from different countries
				usaInfo := &inventoryV1.Part{
					Name:          "USA Part",
					Description:   "Made in USA",
					Price:         2000,
					StockQuantity: 15,
					Category:      inventoryV1.Category_ENGINE,
					Dimensions: &inventoryV1.Dimensions{
						Length: 3, Width: 3, Height: 3, Weight: 30,
					},
					Manufacturer: &inventoryV1.Manufacturer{
						Name: "USA Corp", Country: "USA", Website: "https://usacorp.com",
					},
					Tags: []string{"usa"},
				}
				usaUUID, err := env.InsertTestPartWithData(ctx, usaInfo)
				Expect(err).ToNot(HaveOccurred(), "🟥 Failed to insert test part")

				franceInfo := &inventoryV1.Part{
					Name:          "France Part",
					Description:   "Made in France",
					Price:         2500,
					StockQuantity: 10,
					Category:      inventoryV1.Category_WING,
					Dimensions: &inventoryV1.Dimensions{
						Length: 4, Width: 4, Height: 4, Weight: 40,
					},
					Manufacturer: &inventoryV1.Manufacturer{
						Name: "France Corp", Country: "France", Website: "https://francecorp.com",
					},
					Tags: []string{"france"},
				}
				_, err = env.InsertTestPartWithData(ctx, franceInfo)
				Expect(err).NotTo(HaveOccurred())

				// Act: Filter by USA
				resp, err := invnetoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
					Filter: &inventoryV1.PartsFilter{
						ManufacturerCountries: []string{"USA"},
					},
				})

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.Parts).To(HaveLen(1))
				Expect(resp.Parts[0].Uuid).To(Equal(usaUUID))
				Expect(resp.Parts[0].Manufacturer.Country).To(Equal("USA"))
			})
		})

		Context("when filtering by tags", func() {
			It("should return only parts with matching tags", func() {
				// Arrange: Insert parts with different tags
				wingTagInfo := &inventoryV1.Part{
					Name:          "Wing Part",
					Description:   "Wing with tag",
					Price:         1500,
					StockQuantity: 8,
					Category:      inventoryV1.Category_WING,
					Dimensions: &inventoryV1.Dimensions{
						Length: 2, Width: 2, Height: 2, Weight: 20,
					},
					Manufacturer: &inventoryV1.Manufacturer{
						Name: "Wing Corp", Country: "Germany", Website: "https://wingcorp.com",
					},
					Tags: []string{"wing", "premium"},
				}
				wingUUID, err := env.InsertTestPartWithData(ctx, wingTagInfo)
				Expect(err).NotTo(HaveOccurred())

				engineTagInfo := &inventoryV1.Part{
					Name:          "Engine Part",
					Description:   "Engine with tag",
					Price:         3500,
					StockQuantity: 6,
					Category:      inventoryV1.Category_ENGINE,
					Dimensions: &inventoryV1.Dimensions{
						Length: 3, Width: 3, Height: 3, Weight: 35,
					},
					Manufacturer: &inventoryV1.Manufacturer{
						Name: "Engine Corp", Country: "Japan", Website: "https://enginecorp.com",
					},
					Tags: []string{"engine", "standard"},
				}
				_, err = env.InsertTestPartWithData(ctx, engineTagInfo)
				Expect(err).NotTo(HaveOccurred())

				// Act: Filter by "premium" tag
				resp, err := invnetoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
					Filter: &inventoryV1.PartsFilter{
						Tags: []string{"premium"},
					},
				})

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.Parts).To(HaveLen(1))
				Expect(resp.Parts[0].Uuid).To(Equal(wingUUID))
				Expect(resp.Parts[0].Tags).To(ContainElement("premium"))
			})
		})

		Context("when using multiple filters", func() {
			It("should return parts matching all filter criteria", func() {
				// Arrange: Insert parts with various attributes
				matchingInfo := &inventoryV1.Part{
					Name:          "Premium Wing",
					Description:   "Premium wing from USA",
					Price:         7000,
					StockQuantity: 5,
					Category:      inventoryV1.Category_WING,
					Dimensions: &inventoryV1.Dimensions{
						Length: 8, Width: 4, Height: 2, Weight: 80,
					},
					Manufacturer: &inventoryV1.Manufacturer{
						Name: "USA Wings", Country: "USA", Website: "https://usawings.com",
					},
					Tags: []string{"premium", "certified"},
				}
				matchingUUID, err := env.InsertTestPartWithData(ctx, matchingInfo)
				Expect(err).NotTo(HaveOccurred())

				nonMatchingInfo := &inventoryV1.Part{
					Name:          "Standard Engine",
					Description:   "Standard engine from France",
					Price:         4000,
					StockQuantity: 10,
					Category:      inventoryV1.Category_ENGINE,
					Dimensions: &inventoryV1.Dimensions{
						Length: 5, Width: 3, Height: 3, Weight: 50,
					},
					Manufacturer: &inventoryV1.Manufacturer{
						Name: "France Engines", Country: "France", Website: "https://franceengines.com",
					},
					Tags: []string{"standard"},
				}
				_, err = env.InsertTestPartWithData(ctx, nonMatchingInfo)
				Expect(err).NotTo(HaveOccurred())

				// Act: Filter by category, country, and tag
				resp, err := invnetoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
					Filter: &inventoryV1.PartsFilter{
						Categories:            []inventoryV1.Category{inventoryV1.Category_WING},
						ManufacturerCountries: []string{"USA"},
						Tags:                  []string{"premium"},
					},
				})

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.Parts).To(HaveLen(1))
				Expect(resp.Parts[0].Uuid).To(Equal(matchingUUID))
				Expect(resp.Parts[0].Category).To(Equal(inventoryV1.Category_WING))
				Expect(resp.Parts[0].Manufacturer.Country).To(Equal("USA"))
				Expect(resp.Parts[0].Tags).To(ContainElement("premium"))
			})
		})
	})
})
