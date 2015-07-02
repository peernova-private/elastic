// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// A query that applies a filter to the results of another query.
// For more details, see
// http://www.elasticsearch.org/guide/reference/query-dsl/filtered-query.html
type FilteredQuery struct {
	Query
	query   Query
	filters []Filter
	boost   *float32
}

// Creates a new filtered query.
func NewFilteredQuery(query Query) FilteredQuery {
	q := FilteredQuery{
		query:   query,
		filters: make([]Filter, 0),
	}
	return q
}

func (q FilteredQuery) Filter(filter Filter) FilteredQuery {
	q.filters = append(q.filters, filter)
	return q
}

func (q FilteredQuery) Boost(boost float32) FilteredQuery {
	q.boost = &boost
	return q
}

// Creates the query source for the filtered query.
func (q FilteredQuery) Source() (interface{}, error) {
	// {
	//     "filtered" : {
	//         "query" : {
	//             "term" : { "tag" : "wow" }
	//         },
	//         "filter" : {
	//             "range" : {
	//                 "age" : { "from" : 10, "to" : 20 }
	//             }
	//         }
	//     }
	// }

	source := make(map[string]interface{})

	filtered := make(map[string]interface{})
	source["filtered"] = filtered

	src, err := q.query.Source()
	if err != nil {
		return nil, err
	}
	filtered["query"] = src

	if len(q.filters) == 1 {
		src, err := q.filters[0].Source()
		if err != nil {
			return nil, err
		}
		filtered["filter"] = src
	} else if len(q.filters) > 1 {
		filter := make(map[string]interface{})
		filtered["filter"] = filter
		and := make(map[string]interface{})
		filter["and"] = and
		filters := make([]interface{}, 0)
		for _, f := range q.filters {
			src, err := f.Source()
			if err != nil {
				return nil, err
			}
			filters = append(filters, src)
		}
		and["filters"] = filters
	}

	if q.boost != nil {
		filtered["boost"] = *q.boost
	}

	return source, nil
}
